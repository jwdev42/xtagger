//This file is part of xtagger. ©2023-2026 Jörg Walter.
//This program is free software: you can redistribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.
//
//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with this program.  If not, see <https://www.gnu.org/licenses/>.

package config

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode"
)

// commandParser literals
const (
	litAnd      = "and"
	litAs       = "as"
	litBy       = "by"
	litEOF      = "EOF"
	litFor      = "for"
	litName     = "name"
	litNone     = "none"
	litRecords  = "records"
	litUntagged = "untagged"
)

// commandParserResult holds a parsed command
type commandParserResult struct {
	command      Command    // Application command
	names        []string   // Record names to consider
	paths        []string   // Paths to process
	printRecords bool       // Flag for CommandPrint that controls if records are being printed
	constraints  Constraint // Command constraints
}

// equals returns true if both objects hold equal data.
// This function is only used in the test code.
func (a *commandParserResult) equals(b *commandParserResult) bool {
	if a.command != b.command || a.printRecords != b.printRecords || a.constraints != b.constraints {
		return false
	}
	if !slices.Equal(a.names, b.names) || !slices.Equal(a.paths, b.paths) {
		return false
	}
	return true
}

// CLI command parser.
type commandParser struct {
	tokens []string
	pos    int
	res    *commandParserResult
}

func newCommandParser(args []string) *commandParser {
	return &commandParser{
		tokens: args,
	}
}

// Parser entry point, use this for parsing a command line
func (r *commandParser) start() (*commandParserResult, error) {
	// Reset parser state
	r.pos = 0
	r.res = &commandParserResult{}
	// Start to parse
	command, err := r.parseCommand()
	if err != nil {
		return nil, fmt.Errorf("Command parser failed: %s", err)
	}
	r.res.command = command
	return r.res, nil
}

// Advances to the next token.
func (r *commandParser) adv() {
	r.pos++
}

// Returns the current token and true if it exists for the current position.
// If there is no token left, "EOF" and false will be returned.
func (r *commandParser) tok() (string, bool) {
	if len(r.tokens) > r.pos {
		return r.tokens[r.pos], true
	}
	return litEOF, false
}

func (r *commandParser) error(expected ...string) error {
	tok, _ := r.tok()
	return fmt.Errorf("[Token at index %03d] Expected %q, got %q", r.pos, expected, tok)
}

// Parse a command.
func (r *commandParser) parseCommand() (Command, error) {
	tok, ok := r.tok()
	if !ok {
		return CommandInvalid, r.error("COMMAND")
	}
	command := Command(tok)
	var err error
	switch command {
	case CommandPrint:
		r.adv()
		err = r.parseCommandPrint()
	case CommandTag:
		r.adv()
		err = r.parseCommandTag()
	case CommandUntag:
		r.adv()
		err = r.parseCommandUntag()
	case CommandVerify:
		r.adv()
		err = r.parseCommandVerify()
	case CommandLicenses:
		r.adv()
		err = r.parseCommandLicense()
	default:
		err = fmt.Errorf("Unknown command: %q", command)
	}
	if err != nil {
		return CommandInvalid, err
	}
	return command, nil
}

func (r *commandParser) parseCommandTag() error {
	// Parse "as"
	if err := r.parseLiteral(litAs); err != nil {
		// If "as" is not found, parse tag constraint, then "as"
		r.res.constraints = r.res.constraints.Add(r.parseConstraint(ConstraintUntagged))
		if err := r.parseLiteral(litAs); err != nil {
			return err
		}
	}
	// Parse tag name
	if err := r.parseName(); err != nil {
		return err
	}
	// Parse "for"
	if err := r.parseLiteral(litFor); err != nil {
		return err
	}
	// Parse path(s)
	return r.parsePathsUntilEOF()
}

func (r *commandParser) parseCommandPrint() error {
	// Parse optional constraint
	r.res.constraints = r.res.constraints.Add(r.parseConstraint(ConstraintUntagged))
	// Parse optional literal "records"
	if err := r.parseLiteral(litRecords); err == nil {
		r.res.printRecords = true
	}
	// Parse optional "by" + NAMES
	if err := r.parseLiteral(litBy); err == nil {
		if err := r.parseNames(); err != nil {
			return err
		}
	}
	// Parse "for"
	if err := r.parseLiteral(litFor); err != nil {
		return err
	}
	// Parse PATHS
	return r.parsePathsUntilEOF()
}

func (r *commandParser) parseCommandUntag() error {
	// parse optional names
	if r.isLiteral(litName) {
		if err := r.parseNames(); err != nil {
			return fmt.Errorf("Failed to parse names: %s", err)
		}
	}
	// parse "for"
	if err := r.parseLiteral(litFor); err != nil {
		return err
	}
	// parse PATHS
	return r.parsePathsUntilEOF()
}

func (r *commandParser) parseCommandVerify() error {
	//Parse optional "by" + NAMES
	if err := r.parseLiteral(litBy); err == nil {
		if err := r.parseNames(); err != nil {
			return err
		}
	}
	//Parse "for"
	if err := r.parseLiteral(litFor); err != nil {
		return err
	}
	//Parse PATHS
	return r.parsePathsUntilEOF()
}

func (r *commandParser) parseCommandLicense() error {
	//catch "EOF" token
	_, ok := r.tok()
	if !ok {
		return nil
	}
	return r.error(litEOF)
}

// parseConstraint tries to parse the current token as a constraint.
// If a match is found, parseConstraint will consume the token and will
// return the recognized constraint.
// If no match is found, parseConstraint will not advance the token and
// will return ConstraintNone.
// If the token stream has reached EOF, parseConstraint will return
// ConstraintNone.
// Argument constraints is a variadic slice of all constraints the function
// must detect. If the argument is empty, parseConstraint will not detect
// any constraints.
func (r *commandParser) parseConstraint(constraints ...Constraint) Constraint {
	tok, ok := r.tok()
	if !ok {
		return ConstraintNone
	}
	for _, cs := range constraints {
		if tok == cs.String() {
			r.adv()
			return cs
		}
	}
	return ConstraintNone
}

func (r *commandParser) parsePath() error {
	tok, ok := r.tok()
	if !ok {
		return io.EOF
	}
	if r.res.paths == nil {
		r.res.paths = []string{tok}
	} else {
		r.res.paths = append(r.res.paths, tok)
	}
	r.adv()
	return nil
}

// parseNames parses record names.
func (r *commandParser) parseNames() error {
	if err := r.parseLiteral(litName); err != nil {
		return err
	}
	if err := r.parseName(); err != nil {
		return err
	}
	//parse optional "and"
	if err := r.parseLiteral(litAnd); err != nil {
		//done if token is not "and"
		return nil
	}
	//recurse if optional "and" was parsed
	return r.parseNames()
}

// parseName parses a name expression, this function must only be called
// by parseNames(). Use parseNames() to parse record names!
func (r *commandParser) parseName() error {
	//Closure for name validation
	validateName := func(name string) error {
		if len(name) < 1 {
			return errors.New("Name cannot be empty")
		}
		if strings.TrimSpace(name) != name {
			return errors.New("Name cannot have leading or trailing whitespace")
		}
		for i, ch := range []rune(name) {
			if !unicode.IsPrint(ch) {
				return fmt.Errorf("Character at index %d is not printable", i)
			}
		}
		return nil
	}
	//Check for EOF
	tok, ok := r.tok()
	if !ok {
		return io.EOF
	}
	//Check if token is a valid name
	if err := validateName(tok); err != nil {
		return fmt.Errorf("Invalid name: %s", err)
	}
	//Add name to names slice
	if r.res.names == nil {
		r.res.names = []string{tok}
	} else {
		r.res.names = append(r.res.names, tok)
	}
	r.adv()
	return nil
}

// isLiteral returns true if the current token matches parameter literal.
func (r *commandParser) isLiteral(literal string) bool {
	tok, _ := r.tok()
	return tok == literal
}

// parseLiteral consumes the current token if it matches argument literal,
// then returns nil.
// If the current token does not match literal, an error will be returned
// and the token stream will not be advanced.
func (r *commandParser) parseLiteral(literal string) error {
	tok, ok := r.tok()
	if !ok {
		return io.EOF
	}
	if tok != literal {
		return r.error(literal)
	}
	r.adv()
	return nil
}

// Expects one mandatory path, then parses optional paths until EOF
func (r *commandParser) parsePathsUntilEOF() error {
	//Parse mandatory path
	if err := r.parsePath(); err != nil {
		return err
	}
	//Parse optional paths until EOF
	for {
		if err := r.parsePath(); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

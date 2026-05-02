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
	"strings"
	"unicode"
)

// commandParserResult holds a parsed command
type commandParserResult struct {
	command         Command         // Application command
	paths           []string        // Paths to process
	names           []string        // Record names to consider
	printRecords    bool            // Flag for CommandPrint that controls if records are being printed
	tagConstraint   TagConstraint   // Constraint for command CommandTag
	printConstraint PrintConstraint // Constraint for command CommandPrint
	untagConstraint UntagConstraint // Constraint for command CommandUntag
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
	return "EOF", false
}

func (r *commandParser) error(expected ...string) error {
	tok, _ := r.tok()
	return fmt.Errorf("[Token at index %03d] Expected %q, got %q", r.pos, strings.Join(expected, "\" or \""), tok)
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
	case CommandInvalidate, CommandRevalidate:
		r.adv()
		err = r.parseCommandInvalidateOrRevalidate()
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
	//parse "as"
	if err := r.parseLiteral("as"); err != nil {
		//if "as" is not found, parse tag constraint, then "as"
		if err := r.parseTagConstraint(); err != nil {
			return err
		}
		if err := r.parseLiteral("as"); err != nil {
			return err
		}
	}
	//Parse tag name
	if err := r.parseName(); err != nil {
		return err
	}
	//Parse "for"
	if err := r.parseLiteral("for"); err != nil {
		return err
	}
	//Parse path(s)
	return r.parsePathsUntilEOF()
}

func (r *commandParser) parseCommandPrint() error {
	if err := r.parseLiteral("untagged"); err == nil {
		//Parse "for" after "untagged"
		if err := r.parseLiteral("for"); err != nil {
			return err
		}
		//Parse PATHS
		return r.parsePathsUntilEOF()
	}
	//Parse optional CONSTRAINT, return value can therefore be ignored
	r.parsePrintConstraint()
	//Parse optional literal "records"
	if err := r.parseLiteral("records"); err == nil {
		r.res.printRecords = true
	}
	//Parse optional "by" + NAMES
	if err := r.parseLiteral("by"); err == nil {
		if err := r.parseNames(); err != nil {
			return err
		}
	}
	//Parse "for"
	if err := r.parseLiteral("for"); err != nil {
		return err
	}
	//Parse PATHS
	return r.parsePathsUntilEOF()
}

func (r *commandParser) parseCommandUntag() error {
	//Parse mandatory CONSTRAINT
	if err := r.parseUntagConstraint(); err != nil {
		return err
	}
	//parse "for"
	if err := r.parseLiteral("for"); err != nil {
		return err
	}
	//parse PATHS
	return r.parsePathsUntilEOF()
}

func (r *commandParser) parseCommandInvalidateOrRevalidate() error {
	//parse "all"
	if err := r.parseLiteral("all"); err != nil {
		//parse NAMES if token is not "all"
		if err := r.parseNames(); err != nil {
			return err
		}
	}
	//parse "for"
	if err := r.parseLiteral("for"); err != nil {
		return err
	}
	//parse PATHS
	return r.parsePathsUntilEOF()
}

func (r *commandParser) parseCommandLicense() error {
	//catch "EOF" token
	_, ok := r.tok()
	if !ok {
		return nil
	}
	return r.error(io.EOF.Error())
}

func (r *commandParser) parseTagConstraint() error {
	tok, ok := r.tok()
	if !ok {
		return io.EOF
	}
	switch tok {
	case "untagged":
		r.res.tagConstraint = TagConstraintUntagged
	case "invalid":
		r.res.tagConstraint = TagConstraintInvalid
	default:
		return r.error("untagged", "invalid")
	}
	r.adv()
	return nil
}

func (r *commandParser) parsePrintConstraint() error {
	tok, ok := r.tok()
	if !ok {
		return r.error("invalid", "valid")
	}
	switch tok {
	case "invalid":
		r.res.printConstraint = PrintConstraintInvalid
	case "valid":
		r.res.printConstraint = PrintConstraintValid
	default:
		return r.error("invalid", "valid")
	}
	r.adv()
	return nil
}

func (r *commandParser) parseUntagConstraint() error {
	tok, ok := r.tok()
	if !ok {
		return io.EOF
	}
	switch tok {
	case "all":
		r.res.untagConstraint = UntagConstraintAll
	case "invalid":
		r.res.untagConstraint = UntagConstraintInvalid
	default:
		return r.parseUntagNamesConstraint()
	}
	r.adv()
	return nil
}

func (r *commandParser) parseUntagNamesConstraint() error {
	if err := r.parseNames(); err != nil {
		return err
	}
	//parse optional "if invalid"
	if err := r.parseLiteral("if"); err != nil {
		return nil
	}
	if err := r.parseLiteral("invalid"); err != nil {
		return err
	}
	//Set UntagConstraintInvalid after parsing "if invalid"
	r.res.untagConstraint = UntagConstraintInvalid
	return nil
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

func (r *commandParser) parseNames() error {
	if err := r.parseLiteral("name"); err != nil {
		return err
	}
	if err := r.parseName(); err != nil {
		return err
	}
	//parse optional "and"
	if err := r.parseLiteral("and"); err != nil {
		//done if token is not "and"
		return nil
	}
	//recurse if optional "and" was parsed
	return r.parseNames()
}

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

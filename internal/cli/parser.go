package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// CLI command parser.
type parser struct {
	tokens      []string
	pos         int
	commandLine *CommandLine
}

// Advances to the next token.
func (r *parser) adv() {
	r.pos++
}

// Returns the current token and true if it exists for the current position.
// If there is no token left, "EOF" and false will be returned.
func (r *parser) tok() (string, bool) {
	if len(r.tokens) > r.pos {
		return r.tokens[r.pos], true
	}
	return "EOF", false
}

func (r *parser) error(expected ...string) error {
	tok, _ := r.tok()
	return fmt.Errorf("[Token at index %03d] Expected \"%s\", got \"%s\"", r.pos, strings.Join(expected, "\" or \""), tok)
}

// Parser entry point, parses a command.
func (r *parser) parseCommand() error {
	tok, ok := r.tok()
	if !ok {
		return io.EOF
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
	default:
		err = fmt.Errorf("Unknown command: %q", command)
	}
	if err != nil {
		return err
	}
	r.commandLine.command = command
	return nil
}

func (r *parser) parseCommandTag() error {
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

func (r *parser) parseCommandPrint() error {
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
		r.commandLine.printRecords = true
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

func (r *parser) parseCommandUntag() error {
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

func (r *parser) parseCommandInvalidateOrRevalidate() error {
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

func (r *parser) parseTagConstraint() error {
	tok, ok := r.tok()
	if !ok {
		return io.EOF
	}
	switch tok {
	case "untagged":
		r.commandLine.tagConstraint = TagConstraintUntagged
	case "invalid":
		r.commandLine.tagConstraint = TagConstraintInvalid
	default:
		return r.error("untagged", "invalid")
	}
	r.adv()
	return nil
}

func (r *parser) parsePrintConstraint() error {
	tok, ok := r.tok()
	if !ok {
		return r.error("invalid", "valid")
	}
	switch tok {
	case "invalid":
		r.commandLine.printConstraint = PrintConstraintInvalid
	case "valid":
		r.commandLine.printConstraint = PrintConstraintValid
	default:
		return r.error("invalid", "valid")
	}
	r.adv()
	return nil
}

func (r *parser) parseUntagConstraint() error {
	tok, ok := r.tok()
	if !ok {
		return io.EOF
	}
	switch tok {
	case "all":
		r.commandLine.untagConstraint = UntagConstraintAll
	case "invalid":
		r.commandLine.untagConstraint = UntagConstraintInvalid
	default:
		return r.parseUntagNamesConstraint()
	}
	r.adv()
	return nil
}

func (r *parser) parseUntagNamesConstraint() error {
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
	r.commandLine.untagConstraint = UntagConstraintInvalid
	return nil
}

func (r *parser) parsePath() error {
	tok, ok := r.tok()
	if !ok {
		return io.EOF
	}
	if r.commandLine.paths == nil {
		r.commandLine.paths = []string{tok}
	} else {
		r.commandLine.paths = append(r.commandLine.paths, tok)
	}
	r.adv()
	return nil
}

func (r *parser) parseNames() error {
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

func (r *parser) parseName() error {
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
	if r.commandLine.names == nil {
		r.commandLine.names = []string{tok}
	} else {
		r.commandLine.names = append(r.commandLine.names, tok)
	}
	r.adv()
	return nil
}

func (r *parser) parseLiteral(literal string) error {
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
func (r *parser) parsePathsUntilEOF() error {
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

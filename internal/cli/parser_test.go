package cli

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	var tests = map[*[]string]*CommandLine{
		{"tag", "as", "foo", "for", "/tmp"}: {
			command:       CommandTag,
			names:         []string{"foo"},
			paths:         []string{"/tmp"},
			tagConstraint: TagConstraintNone,
		},
		{"tag", "as", "foo", "for", "tmp", "tmp2"}: {
			command:       CommandTag,
			names:         []string{"foo"},
			paths:         []string{"tmp", "tmp2"},
			tagConstraint: TagConstraintNone,
		},
		{"tag", "as", "foo", "for", "tmp", "tmp2", "tmp3"}: {
			command:       CommandTag,
			names:         []string{"foo"},
			paths:         []string{"tmp", "tmp2", "tmp3"},
			tagConstraint: TagConstraintNone,
		},
		{"tag", "untagged", "as", "foo", "for", "tmp"}: {
			command:       CommandTag,
			names:         []string{"foo"},
			paths:         []string{"tmp"},
			tagConstraint: TagConstraintUntagged,
		},
		{"tag", "invalid", "as", "foo", "for", "/tmp"}: {
			command:       CommandTag,
			names:         []string{"foo"},
			paths:         []string{"/tmp"},
			tagConstraint: TagConstraintInvalid,
		},
		{"untag", "all", "for", "/tmp"}: {
			command:         CommandUntag,
			names:           nil,
			paths:           []string{"/tmp"},
			untagConstraint: UntagConstraintAll,
		},
		{"untag", "all", "for", "tmp", "tmp2"}: {
			command:         CommandUntag,
			names:           nil,
			paths:           []string{"tmp", "tmp2"},
			untagConstraint: UntagConstraintAll,
		},
		{"untag", "invalid", "for", "tmp", "tmp2"}: {
			command:         CommandUntag,
			names:           nil,
			paths:           []string{"tmp", "tmp2"},
			untagConstraint: UntagConstraintInvalid,
		},
		{"untag", "name", "example", "for", "tmp", "tmp2"}: {
			command:         CommandUntag,
			names:           []string{"example"},
			paths:           []string{"tmp", "tmp2"},
			untagConstraint: UntagConstraintNone,
		},
		{"untag", "name", "example", "and", "name", "foobar", "for", "tmp", "tmp2"}: {
			command:         CommandUntag,
			names:           []string{"example", "foobar"},
			paths:           []string{"tmp", "tmp2"},
			untagConstraint: UntagConstraintNone,
		},
		{"untag", "name", "example", "if", "invalid", "for", "tmp", "tmp2"}: {
			command:         CommandUntag,
			names:           []string{"example"},
			paths:           []string{"tmp", "tmp2"},
			untagConstraint: UntagConstraintInvalid,
		},
		{"print", "for", "/tmp"}: {
			command:         CommandPrint,
			names:           nil,
			paths:           []string{"/tmp"},
			printConstraint: PrintConstraintNone,
		},
		{"print", "for", "for", "test"}: {
			command:         CommandPrint,
			names:           nil,
			paths:           []string{"for", "test"},
			printConstraint: PrintConstraintNone,
		},
		{"print", "valid", "for", "test"}: {
			command:         CommandPrint,
			names:           nil,
			paths:           []string{"test"},
			printConstraint: PrintConstraintValid,
		},
		{"print", "invalid", "for", "test"}: {
			command:         CommandPrint,
			names:           nil,
			paths:           []string{"test"},
			printConstraint: PrintConstraintInvalid,
		},
		{"print", "untagged", "for", "test"}: {
			command:         CommandPrint,
			names:           nil,
			paths:           []string{"test"},
			printConstraint: PrintConstraintUntagged,
		},
		{"invalidate", "all", "for", "test"}: {
			command: CommandInvalidate,
			names:   nil,
			paths:   []string{"test"},
		},
		{"invalidate", "name", "foo", "for", "test"}: {
			command: CommandInvalidate,
			names:   []string{"foo"},
			paths:   []string{"test"},
		},
		{"invalidate", "name", "foo", "and", "name", "bar", "for", "test"}: {
			command: CommandInvalidate,
			names:   []string{"foo", "bar"},
			paths:   []string{"test"},
		},
		{"revalidate", "all", "for", "test"}: {
			command: CommandRevalidate,
			names:   nil,
			paths:   []string{"test"},
		},
		{"revalidate", "name", "foo", "for", "test"}: {
			command: CommandRevalidate,
			names:   []string{"foo"},
			paths:   []string{"test"},
		},
		{"revalidate", "name", "foo", "and", "name", "bar", "for", "test"}: {
			command: CommandRevalidate,
			names:   []string{"foo", "bar"},
			paths:   []string{"test"},
		},
	}
	for tokens, blueprint := range tests {
		var p = &parser{
			tokens:      *tokens,
			commandLine: new(CommandLine),
		}
		if err := p.parseCommand(); err != nil {
			t.Errorf("Parser error for command \"%s\": %s", strings.Join(*tokens, " "), err)
			continue
		}
		if err := blueprint.mustEqual(p.commandLine); err != nil {
			t.Errorf("Blueprint doesn't match result for command \"%s\": %s", strings.Join(*tokens, " "), err)
		}
	}
}

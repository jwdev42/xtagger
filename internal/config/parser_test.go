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
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	var tests = map[*[]string]*commandParserResult{
		{"tag", "as", "foo", "for", "/tmp"}: {
			command:     CommandTag,
			names:       []string{"foo"},
			paths:       []string{"/tmp"},
			constraints: ConstraintNone,
		},
		{"tag", "as", "foo", "for", "tmp", "tmp2"}: {
			command:     CommandTag,
			names:       []string{"foo"},
			paths:       []string{"tmp", "tmp2"},
			constraints: ConstraintNone,
		},
		{"tag", "as", "foo", "for", "tmp", "tmp2", "tmp3"}: {
			command:     CommandTag,
			names:       []string{"foo"},
			paths:       []string{"tmp", "tmp2", "tmp3"},
			constraints: ConstraintNone,
		},
		{"tag", "untagged", "as", "foo", "for", "tmp"}: {
			command:     CommandTag,
			names:       []string{"foo"},
			paths:       []string{"tmp"},
			constraints: ConstraintUntagged,
		},
		{"untag", "for", "/tmp"}: {
			command:     CommandUntag,
			names:       nil,
			paths:       []string{"/tmp"},
			constraints: ConstraintNone,
		},
		{"untag", "for", "tmp", "tmp2"}: {
			command:     CommandUntag,
			names:       nil,
			paths:       []string{"tmp", "tmp2"},
			constraints: ConstraintNone,
		},
		{"untag", "name", "example", "for", "tmp", "tmp2"}: {
			command:     CommandUntag,
			names:       []string{"example"},
			paths:       []string{"tmp", "tmp2"},
			constraints: ConstraintNone,
		},
		{"untag", "name", "example", "and", "name", "foobar", "for", "tmp", "tmp2"}: {
			command:     CommandUntag,
			names:       []string{"example", "foobar"},
			paths:       []string{"tmp", "tmp2"},
			constraints: ConstraintNone,
		},
		{"print", "for", "/tmp"}: {
			command:     CommandPrint,
			names:       nil,
			paths:       []string{"/tmp"},
			constraints: ConstraintNone,
		},
		{"print", "for", "for", "test"}: {
			command:     CommandPrint,
			names:       nil,
			paths:       []string{"for", "test"},
			constraints: ConstraintNone,
		},
		{"print", "untagged", "for", "test"}: {
			command:     CommandPrint,
			names:       nil,
			paths:       []string{"test"},
			constraints: ConstraintUntagged,
		},
	}
	for tokens, blueprint := range tests {
		p := newCommandParser(*tokens)
		res, err := p.start()
		if err != nil {
			t.Errorf("Parser error for command \"%s\": %s", strings.Join(*tokens, " "), err)
			continue
		}
		if !blueprint.equals(res) {
			t.Errorf("Blueprint doesn't match result for command %q", strings.Join(*tokens, " "))
		}
	}
}

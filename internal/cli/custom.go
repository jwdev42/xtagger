package cli

import (
	"errors"
	"fmt"
	"unicode"
)

// Used to parse the identifier for an xtag
type name string

func (r *name) String() string {
	if r == nil {
		return ""
	}
	return string(*r)
}

func (r *name) Get() any {
	if r == nil {
		return nil
	}
	return string(*r)
}

func (r *name) Set(s string) error {
	if len(s) < 1 {
		return errors.New("name cannot be empty")
	}
	for i, ch := range []rune(s) {
		if !unicode.IsPrint(ch) {
			return fmt.Errorf("Invalid character at index %d", i)
		}
	}
	*r = name(s)
	return nil
}

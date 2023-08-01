package cli

import(
	"errors"
	"unicode"
)

type name string

func (r *name) String() string {
	if r == nil {
		return ""
	}
	return *r
}

func (r *name) Set(s string) error {
	if len(s) < 1 {
		return errors.New("name cannot be empty")
	}
	for i, ch := range []rune(s) {
		if !unicode.IsPrint(r) {
			return fmt.Errorf("Invalid character at index %d", i)
		}
	}
	r* = name(s)
	return nil
}

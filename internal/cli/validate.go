package cli

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func validateName(name string) error {
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

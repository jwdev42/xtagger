package cli

import (
	"errors"
	"fmt"
	"github.com/jwdev42/xtagger/internal/hashes"
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

type hash hashes.Algo

func (r *hash) String() string {
	if r == nil {
		return ""
	}
	return string(*r)
}

func (r *hash) Get() any {
	if r == nil {
		return nil
	}
	return hashes.Algo(*r)
}

func (r *hash) Set(s string) error {
	if len(s) < 1 {
		return errors.New("hash cannot be empty")
	}
	algo, err := hashes.ParseAlgo(s)
	if err != nil {
		return err
	}
	*r = hash(algo)
	return nil
}

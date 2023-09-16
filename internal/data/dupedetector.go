package data

import (
	"errors"
	"fmt"
	"hash"
	"io"
)

var DupeDetected = errors.New("Dupe detected")

type DupeDetector map[string]bool

func (r DupeDetector) Register(stream io.Reader, hash hash.Hash) error {
	hash.Reset()
	if _, err := io.Copy(hash, stream); err != nil {
		return err
	}
	sum := fmt.Sprintf("%x", hash.Sum(nil))
	if r[sum] {
		return DupeDetected
	}
	r[sum] = true
	return nil
}

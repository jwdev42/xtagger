package hashes

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
	"hash"
)

const (
	SHA256    Algo = "SHA256"    //SHA 256
	RIPEMD160      = "RIPEMD160" //RIPEMD 160
	SHA3256        = "SHA3256"   //SHA3 256
)

// Represents a name of a particular hashing algorithm at a particluar length.
type Algo string

// Returns the Algo corresponding to Name, returns an error if Name
// does not represent an existing Algo
func ParseAlgo(Name string) (Algo, error) {
	switch Name {
	case "SHA256", "sha256", "SHA_256", "sha_256":
		return SHA256, nil
	case "RIPEMD160", "ripemd160", "ripemd_160", "RIPEMD_160":
		return RIPEMD160, nil
	case "SHA3256", "sha3256", "SHA3_256", "sha3_256":
		return SHA3256, nil
	}
	return "", fmt.Errorf("Unknown hashing algorithm \"%s\"", Name)
}

// Returns a usable hash.Hash interface. If the receiver is not a valid name for
// a supported hash function, the method panics.
func (r Algo) New() hash.Hash {
	switch r {
	case SHA256:
		return sha256.New()
	case RIPEMD160:
		return ripemd160.New()
	case SHA3256:
		return sha3.New256()
	}
	panic(fmt.Errorf("Receiver has an invalid value: \"%s\"", r))
}

func (r *Algo) UnmarshalText(text []byte) error {
	if text == nil {
		return errors.New("Algo cannot be nil")
	}
	algo, err := ParseAlgo(string(text))
	if err != nil {
		return err
	}
	*r = algo
	return nil
}

func (r Algo) Validate() error {
	switch r {
	case SHA256, RIPEMD160, SHA3256:
		return nil
	}
	return fmt.Errorf("Invalid hashing algorithm \"%s\"", r)
}

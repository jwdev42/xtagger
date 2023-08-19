package io

import (
	"hash"
	"io"
)

func Hash(src io.Reader, hash hash.Hash) error {
	buf := make([]byte, bufsize)
	for true {
		r, err := src.Read(buf)
		if r > 0 {
			//write to hash
			hash.Write(buf[:r])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

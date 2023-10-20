package hashes

import (
	"github.com/jwdev42/xtagger/internal/global"
	"hash"
	"io"
)

func Hash(src io.Reader, hash hash.Hash) error {
	buf := make([]byte, global.BufSize)
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

func MultiHash(src io.Reader, hashMap map[Algo]hash.Hash) error {
	buf := make([]byte, global.BufSize)
	var n int
	var readErr error
	for {
		n, readErr = src.Read(buf)
		if n > 0 {
			for _, hasher := range hashMap {
				_, err := hasher.Write(buf[:n])
				if err != nil {
					return err
				}
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}
	return nil
}

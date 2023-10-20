package hashes

import (
	"github.com/jwdev42/xtagger/internal/global"
	"hash"
	"io"
)

func HashCopy(dst io.Writer, src io.Reader, hash hash.Hash) (written int64, err error) {
	buf := make([]byte, global.BufSize)
	for true {
		r, err := src.Read(buf)
		if r > 0 {
			//write to dest
			w, err := dst.Write(buf[:r])
			written = written + int64(w)
			if err != nil {
				return written, err
			}
			//write to hash
			hash.Write(buf[:r])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return written, err
		}
	}
	return written, nil
}

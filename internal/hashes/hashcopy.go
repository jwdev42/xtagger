//This file is part of xtagger. ©2023 Jörg Walter.
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

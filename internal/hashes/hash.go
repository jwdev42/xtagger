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

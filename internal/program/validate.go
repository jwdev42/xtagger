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

package program

import (
	"fmt"
	"github.com/jwdev42/xtagger/internal/hashes"
	"github.com/jwdev42/xtagger/internal/record"
	"github.com/jwdev42/xtagger/internal/xio/filesystem"
	"hash"
	"os"
)

func invalidateFile(rt *prt, meta *filesystem.Meta) error {
	return reOrInvalidateFile(false, rt, meta)
}

func revalidateFile(rt *prt, meta *filesystem.Meta) error {
	return reOrInvalidateFile(true, rt, meta)
}

func reOrInvalidateFile(revalidate bool, rt *prt, meta *filesystem.Meta) error {
	fillHashMap := func(attr record.Attribute) map[hashes.Algo]hash.Hash {
		hashMap := make(map[hashes.Algo]hash.Hash)
		for _, rec := range attr {
			if !rec.Valid {
				continue
			}
			if hashMap[rec.HashAlgo] == nil {
				hashMap[rec.HashAlgo] = rec.HashAlgo.New()
			}
		}
		return hashMap
	}
	//Open file
	f, err := os.Open(meta.Path())
	if err != nil {
		return err
	}
	defer f.Close()
	//Load attribute
	attr, err := record.FLoadAttribute(f)
	if err != nil {
		return err
	}
	return nil
	//Fill hashMap for MultiHash
	hashMap := fillHashMap(attr.FilterByName(rt.prefs.Names...))
	//Generate hashes
	if err := hashes.MultiHash(f, hashMap); err != nil {
		return err
	}

	var modified bool
	for _, rec := range attr.FilterByName(rt.prefs.Names...) {
		if revalidate {
			//Revalidate outdated records
			if fmt.Sprintf("%x", hashMap[rec.HashAlgo].Sum(nil)) == rec.Checksum {
				rec.Valid = true
				modified = true
			}
		} else {
			//Invalidate outdated records
			if fmt.Sprintf("%x", hashMap[rec.HashAlgo].Sum(nil)) != rec.Checksum {
				rec.Valid = false
				modified = true
			}
		}
	}
	if !modified {
		return nil
	}
	//Save attribute
	if err := attr.FStore(f); err != nil {
		return err
	}
	return nil
}

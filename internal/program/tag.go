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
	"github.com/jwdev42/logger"
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/hashes"
	"github.com/jwdev42/xtagger/internal/record"
	"github.com/jwdev42/xtagger/internal/softerrors"
	"io/fs"
	"os"
	"path/filepath"
)

func tagFile(parent string, info fs.FileInfo) error {
	path := filepath.Join(parent, info.Name())
	name := commandLine.Names()[0]
	algo := commandLine.FlagHash()
	constraint := commandLine.TagConstraint()
	//Open file
	f, err := os.Open(path)
	if err != nil {
		return softerrors.Consume(err)
	}
	defer f.Close()
	//Load attribute
	attr, err := record.FLoadAttribute(f)
	if err != nil {
		return softerrors.Consume(err)
	}
	//Process untagged constraint
	if constraint == cli.TagConstraintUntagged && len(attr) > 0 {
		return nil //Skip already tagged files
	}
	//Process invalid constraint
	if constraint == cli.TagConstraintInvalid {
		for _, rec := range attr {
			if rec.Valid {
				//Skip files that have a valid record
				return nil
			}
		}
	}
	//Check if a record with the designated name already exists
	if attr.Exists(name) {
		return softerrors.Consume(fmt.Errorf("Record \"%s\" already exists for path \"%s\"", name, path))
	}
	//Hash file
	logger.Default().Infof("Hashing file %s", path)
	hash := algo.New()
	if err := hashes.Hash(f, hash); err != nil {
		return softerrors.Consume(err)
	}
	logger.Default().Infof("Successfully hashed file %s", path)
	//Create record
	rec := record.NewRecord()
	rec.Checksum = fmt.Sprintf("%x", hash.Sum(nil))
	rec.HashAlgo = algo
	rec.Valid = true
	//Add record to attribute
	attr[name] = rec
	//Save attribute
	if err := attr.FStore(f); err != nil {
		return softerrors.Consume(err)
	}
	//Print path if print0 is active
	if commandLine.FlagPrint0() {
		if _, err := printMe.Print0(path); err != nil {
			return softerrors.Consume(err)
		}
	}
	return nil
}

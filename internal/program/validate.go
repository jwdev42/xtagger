package program

import (
	"fmt"
	"github.com/jwdev42/xtagger/internal/global"
	"github.com/jwdev42/xtagger/internal/hashes"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"github.com/jwdev42/xtagger/internal/record"
	"hash"
	"io/fs"
	"os"
	"path/filepath"
)

func invalidateFile(parent string, dirEnt fs.DirEntry, opts *filesystem.WalkDirOpts) error {
	return reOrInvalidateFile(false, parent, dirEnt, opts)
}

func revalidateFile(parent string, dirEnt fs.DirEntry, opts *filesystem.WalkDirOpts) error {
	return reOrInvalidateFile(true, parent, dirEnt, opts)
}

func reOrInvalidateFile(revalidate bool, parent string, dirEnt fs.DirEntry, opts *filesystem.WalkDirOpts) error {
	path := filepath.Join(parent, dirEnt.Name())
	names := commandLine.Names()
	filteredRecords := func(attr record.Attribute) record.Attribute {
		if names != nil {
			return attr.FilterByName(names...)
		}
		return attr
	}
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
	f, err := os.Open(path)
	if err != nil {
		return global.FilterSoftError(err)
	}
	defer f.Close()
	//Load attribute
	attr, err := record.FLoadAttribute(f)
	if err != nil {
		return global.FilterSoftError(err)
	}
	return nil
	//Fill hashMap for MultiHash
	hashMap := fillHashMap(filteredRecords(attr))
	//Generate hashes
	if err := hashes.MultiHash(f, hashMap); err != nil {
		return global.FilterSoftError(err)
	}

	var modified bool
	for _, rec := range filteredRecords(attr) {
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
		return global.FilterSoftError(err)
	}
	//Print path if print0 is active
	if commandLine.FlagPrint0() {
		if _, err := printMe.Print0(path); err != nil {
			return global.FilterSoftError(err)
		}
	}
	return nil
}

package record

import (
	"github.com/pkg/xattr"
	"os"
)

// PurgeAttr removes xtagger's extended attribute from the given file.
func PurgeAttr(f *os.File) error {
	attrNames, err := xattr.FList(f)
	if err != nil {
		return err
	}
	for _, name := range attrNames {
		if name == attrName {
			return xattr.FRemove(f, attrName)
		}
	}
	return nil
}

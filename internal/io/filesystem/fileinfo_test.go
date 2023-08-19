package filesystem

import (
	"testing"
)

func TestStat(t *testing.T) {
	const path = "test"
	info, _ := Stat(path)
	if err := info.Error(); err != nil {
		t.Errorf("Stat returned an error: %s", err)
	}
	t.Logf("FileInfo for path \"%s\": %+v", path, *info)
}

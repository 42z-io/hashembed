package hashembed

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDirEntry struct{}

func (m MockDirEntry) Name() string {
	return "name.ext"
}

func (m MockDirEntry) RootPath() string {
	return "root"
}

func (m MockDirEntry) IsDir() bool {
	return false
}

func (m MockDirEntry) Type() fs.FileMode {
	return 0
}

func (m MockDirEntry) NameAndExtension() (string, string) {
	return "name", "ext"
}

func (m MockDirEntry) Info() (fs.FileInfo, error) {
	return nil, errors.New("unimplemented")
}

func TestExtensionRenamer(t *testing.T) {
	var entry MockDirEntry
	assert.Equal(
		t,
		"root/name.hash.ext",
		ExtensionRenamer(
			PathedDirEntry{
				entry:    entry,
				rootPath: "root",
			},
			"hash",
		),
		"filename should be renamed by adding the hash before the extension",
	)
}

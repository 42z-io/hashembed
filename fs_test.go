package hashembed

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/*
var testEmbed embed.FS

var testHashEmbed *FS
var testHashEmbedOptions *FS

func init() {
	testHashEmbed, _ = Generate(testEmbed)
	testHashEmbedOptions, _ = Generate(testEmbed, Config{
		Hasher:  Crc32Hasher,
		Renamer: FullNameRenamer,
	})
}

func TestReverse(t *testing.T) {
	assert.Equal(
		t,
		"testdata/7f2cded6.css",
		testHashEmbedOptions.Reverse("testdata/test.css"),
		"reverse path should be calculated correctly",
	)
	assert.Equal(
		t,
		"testdata/test.txt",
		testHashEmbedOptions.Reverse("testdata/test.txt"),
		"reverse path should remain unchanged when file not hashed",
	)
}

func TestForward(t *testing.T) {
	assert.Equal(
		t,
		"testdata/test.css",
		testHashEmbedOptions.Forward("testdata/7f2cded6.css"),
		"forward path should be calculated correctly",
	)
	assert.Equal(
		t,
		"testdata/test.txt",
		testHashEmbedOptions.Forward("testdata/test.txt"),
		"forward path should remain unchanged when file not hashed",
	)
}

func TestReadFile(t *testing.T) {
	data, err := testHashEmbedOptions.ReadFile("testdata/7f2cded6.css")
	assert.Equal(
		t,
		nil,
		err,
		"no error should occur reading the file",
	)
	assert.Equal(
		t,
		"body { width: 100%; }\n",
		string(data[:]),
		"file should have correct data",
	)
}

func TestReadFileNested(t *testing.T) {
	data, err := testHashEmbed.ReadFile("testdata/folder/test2_123@#%(!.css")
	assert.Equal(
		t,
		nil,
		err,
		"no error should occur reading the nested file",
	)
	assert.Equal(
		t,
		"body { width: 5000px; }\n",
		string(data[:]),
		"nested file should have correct data",
	)
}

func TestReadDir(t *testing.T) {
	data, err := testHashEmbed.ReadDir("testdata")
	assert.Equal(
		t,
		nil,
		err,
		"no error should occur reading a directory",
	)
	assert.Equal(
		t,
		"folder",
		data[0].Name(),
		"folder should be in directory",
	)
	assert.Equal(
		t,
		true,
		data[0].IsDir(),
		"folder should be a directory",
	)
	assert.Equal(
		t,
		"test.css",
		data[1].Name(),
		"test.css should be in directory",
	)
	assert.Equal(
		t,
		"test.txt",
		data[2].Name(),
		"test.txt should be in directory",
	)
}

func TestOpenFile(t *testing.T) {
	f, err := testHashEmbed.Open("testdata/test.txt")
	assert.Equal(
		t,
		nil,
		err,
		"no error should occur opening a file",
	)

	stat, err := f.Stat()
	assert.Equal(
		t,
		nil,
		err,
		"no error should occur when get the file info for an opened file",
	)

	assert.Equal(
		t,
		"test.txt",
		stat.Name(),
		"text.txt should be the name of the opened file",
	)
}

func BenchmarkGenerate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Generate(testEmbed)
	}
}

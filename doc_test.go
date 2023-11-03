package hashembed

import (
	"embed"
	"fmt"
)

//go:embed testdata/*
var embedded embed.FS

var hashedEmbeded *HashedFS

func init() {
	hashedEmbeded, _ = Generate(embedded)
}

func Example() {
	// use go:embed
	// var embeded embed.FS
	embedded, _ := Generate(embedded)
	path := embedded.GetHashedPath("testdata/test.css")
	data, _ := embedded.ReadFile(path)
	fmt.Printf("%s\n%s\n", path, string(data[:]))
	// Output: testdata/test.8d77f04c3be2abcd554f262130ba6c30f277318e66588b6a0d95f476c4ae7c48.css
	// body { width: 100%; }
}

func Example_configured() {
	// use go:embed
	// var embeded embed.FS
	embedded, _ := Generate(embedded, Config{
		// Extensions not in this list will not be given content-hashes
		AllowedExtensions: []string{"css", "txt"},
		// Mechanism to control the hash
		Hasher: Crc32Hasher{},
		// Mechanism to control the naming of the content-hashed files
		Renamer: FullNameRenamer,
	})
	path := embedded.GetHashedPath("testdata/test.css")
	data, _ := embedded.ReadFile(path)
	fmt.Printf("%s\n%s\n", path, string(data[:]))
	// Output: testdata/7f2cded6.css
	// body { width: 100%; }
}

func ExampleHashedFS_GetActualPath() {
	fmt.Println(
		hashedEmbeded.GetActualPath("testdata/test.8d77f04c3be2abcd554f262130ba6c30f277318e66588b6a0d95f476c4ae7c48.css"),
	)
	// Output: testdata/test.css
}

func ExampleHashedFS_GetHashedPath() {
	fmt.Println(
		hashedEmbeded.GetHashedPath("testdata/test.css"),
	)
	// Output: testdata/test.8d77f04c3be2abcd554f262130ba6c30f277318e66588b6a0d95f476c4ae7c48.css
}

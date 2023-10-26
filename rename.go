package hashembed

import "fmt"

// Functional interface for a custom FileRenamer
type FileRenamer func(entry PathedDirEntry, hash string) string

// Rename a file by injecting the hash before the extension.
//
//	my/path/test.css -> my/path/test.$hash.css
func ExtensionRenamer(file PathedDirEntry, hash string) string {
	name, ext := file.NameAndExtension()
	return fmt.Sprintf("%s/%s.%s.%s", file.RootPath(), name, hash, ext)
}

// Rename a file by replacing the name with the hash
//
//	my/path/test.css -> my/path/$hash.css
func FullNameRenamer(file PathedDirEntry, hash string) string {
	_, ext := file.NameAndExtension()
	return fmt.Sprintf("%s/%s.%s", file.RootPath(), hash, ext)
}

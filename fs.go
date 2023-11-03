// [hashembed] is an [embed.FS] with support for reading files with virtual content hashes embedded in the file name.
//
// [hashembed] is useful if you are embedding static assets directly into your application and want to
// facilitate serving these files with very long duration client-side caching.
//
// # File Hashing
//
// Files are hashed when you call [Generate].
//
// You can provide a custom file hasher by providing a function that matches [FileHasher].
//
// There are several built-in hashers:
//   - [Sha256Hasher] (default)
//   - [Crc32Hasher]
//
// # File Renaming
//
// Files are renamed to include their hash when you call [Generate].
//
// You can provide a custom file renamer by providing a function that matches [FileRenamer].
//
// There are two built-in renaming mechanims:
//   - [ExtensionRenamer] (default)
//   - [FullNameRenamer]
//
// # Examples
//
// [embed]: https://pkg.go.dev/embed
package hashembed

import (
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"io/fs"
	"slices"
)

// HashedFS is an [embed.FS] with support for reading files with virtual content hashes embedded in the file name.
type HashedFS struct {
	fs               embed.FS          // underlying embed.FS
	actualPathLookup map[string]string // lookups for the hashed path => actual path
	hashedPathLookup map[string]string // lookups for the actual path => hashed path
	integrityLookup  map[string]string // lookups for the actual path => integrity hash (sha-256)
	cfg              Config            // configuration options for the hashed embed
}

// Get the integrity (subresource integrity) for a file as base64.
func (f HashedFS) getIntegrityBase64(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))

}

// Initialize a file by generating a hash, renaming (aliasing), and adding it to the lookup.
func (f HashedFS) initializeFile(file PathedDirEntry) error {
	_, ext := file.NameAndExtension()
	if !slices.Contains(f.cfg.AllowedExtensions, ext) {
		return nil
	}

	data, err := f.fs.ReadFile(file.FullPath())
	if err != nil {
		return err
	}

	hash, err := f.cfg.Hasher(data)
	if err != nil {
		return err
	}

	hashedPath := f.cfg.Renamer(file, hash)
	if err != nil {
		return err
	}

	fullPath := file.FullPath()
	f.actualPathLookup[hashedPath] = fullPath
	f.hashedPathLookup[fullPath] = hashedPath
	f.integrityLookup[fullPath] = f.getIntegrityBase64(data)
	return nil
}

// Initialize a path (could be file or directory) within the embed.FS.
func (f HashedFS) initializePath(root PathedDirEntry) error {
	rootPath := root.FullPath()
	entries, err := f.fs.ReadDir(rootPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		pathEntry := NewPathedDirEntry(entry, rootPath)
		if !entry.IsDir() {
			if err := f.initializeFile(pathEntry); err != nil {
				return err
			}
		} else {
			if err := f.initializePath(pathEntry); err != nil {
				return err
			}
		}
	}
	return nil
}

// Initialize the [HashedFS] by iterating over the files in the embed.FS.
func (f HashedFS) initialize() error {
	entries, err := f.fs.ReadDir(".")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		pathEntry := NewPathedDirEntry(entry, "")
		if err := f.initializePath(pathEntry); err != nil {
			return err
		}
	}
	return nil
}

// Generate will create a new instance of [HashedFS] using [Config] (if provided) or [ConfigDefault] if not provided.
func Generate(fs embed.FS, cfgs ...Config) (*HashedFS, error) {
	cfg := ConfigDefault
	if len(cfgs) > 0 {
		cfg = cfgs[0]

		if cfg.Hasher == nil {
			cfg.Hasher = ConfigDefault.Hasher
		}

		if cfg.Renamer == nil {
			cfg.Renamer = ConfigDefault.Renamer
		}

		if cfg.AllowedExtensions == nil {
			cfg.AllowedExtensions = ConfigDefault.AllowedExtensions
		}
	}

	hashedEmbed := &HashedFS{
		fs:               fs,
		actualPathLookup: make(map[string]string),
		hashedPathLookup: make(map[string]string),
		integrityLookup:  make(map[string]string),
		cfg:              cfg,
	}

	hashedEmbed.initialize()
	return hashedEmbed, nil
}

// GetActualPath will convert the content hashed path into the actual path.
//
// If the actual path is not found it will return the provided path.
func (f HashedFS) GetActualPath(path string) string {
	if lookup, ok := f.actualPathLookup[path]; ok {
		return lookup
	}
	return path
}

// GetHashedPath will convert the actual path into the content hashed path.
//
// If the hashed path is not found it will return the provided path.
func (f HashedFS) GetHashedPath(path string) string {
	if lookup, ok := f.hashedPathLookup[path]; ok {
		return lookup
	}
	return path
}

// GetIntegrity will get the SHA-256 integrity hash (base64 encoded) for the specified path.
//
// Will only find files matched by the [Config.AllowedExtensions] list.
//
// If the hashed path is not found it will return a blank string.
func (f HashedFS) GetIntegrity(path string) string {
	if lookup, ok := f.integrityLookup[path]; ok {
		return lookup
	}
	return ""
}

// See [embed.FS.Open]
//
// This will call [GetActualPath] on the file to get the correct name.
func (f HashedFS) Open(name string) (fs.File, error) {
	return f.fs.Open(f.GetActualPath(name))
}

// See [embed.FS.ReadDir]
//
// Note: This will only return files that actually exist in the [embed.FS] - hashed files are "virtual"
func (f HashedFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return f.fs.ReadDir(name)
}

// See [embed.FS]
//
// This will call [HashedFS.GetActualPath] on the file to get the correct name.
func (f HashedFS) ReadFile(name string) ([]byte, error) {
	return f.fs.ReadFile(f.GetActualPath(name))
}

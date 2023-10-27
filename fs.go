// [hashembed] is a thin wrapper around [embed] to allow accessing files with a content hash.
//
// [hashembed] is useful if you are embedding static assets directly into your application and want to
// facilitate serving these files with very long duration client-side caching.
//
// # File Hashing
//
// Files are hashed when you call [FS.Generate].
//
// You can provide a custom file hasher by providing a function that matches [FileHasher].
//
// There are several built-in hashers:
//   - [Sha256Hasher] (default)
//   - [Crc32Hasher]
//
// # File Renaming
//
// Files are renamed to include their hash when you call [FS.Generate].
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
	"embed"
	"io/fs"
	"slices"
)

// FS is the "emulated" embed.FS interface with support for content hashes
type FS struct {
	fs               embed.FS          // underlying embed.FS
	actualPathLookup map[string]string // lookups for the hashed path => actual path
	hashedPathLookup map[string]string // lookups for the actual path => hashed path
	cfg              Config            // configuration options for the hashed embed
}

// Initialize a file by generating a hash, renaming (aliasing), and adding it to the lookup.
func (f FS) initializeFile(file PathedDirEntry) error {
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
	return nil
}

// Initialize a path (could be file or directory) within the embed.FS.
func (f FS) initializePath(root PathedDirEntry) error {
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

// Initialize the FS by iterating over the files in the embed.FS.
func (f FS) initialize() error {
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

// Generate will create a new instance of [FS] using [Config] (if provided) or [ConfigDefault] if not provided.
func Generate(fs embed.FS, cfgs ...Config) (*FS, error) {
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

	hashedEmbed := &FS{
		fs:               fs,
		actualPathLookup: make(map[string]string),
		hashedPathLookup: make(map[string]string),
		cfg:              cfg,
	}

	hashedEmbed.initialize()
	return hashedEmbed, nil
}

// GetActualPath will convert the content hashed path into the actual path.
func (f FS) GetActualPath(path string) string {
	if lookup, ok := f.actualPathLookup[path]; ok {
		return lookup
	}
	return path
}

// GetHashedPath will convert the actual path into the content hashed path.
func (f FS) GetHashedPath(path string) string {
	if lookup, ok := f.hashedPathLookup[path]; ok {
		return lookup
	}
	return path
}

// Wrapper for embed.FS.Open.
//
// This will call [GetActualPath] on the file to get the correct name.
func (f FS) Open(name string) (fs.File, error) {
	return f.fs.Open(f.GetActualPath(name))
}

// Wrapper for embed.FS.ReadDir.
//
// This will NOT [GetActualPath] the name - ReadDir is not currently supported by [hashembed].
func (f FS) ReadDir(name string) ([]fs.DirEntry, error) {
	return f.fs.ReadDir(name)
}

// Wrapper for embed.FS.ReadFile.
//
// This will call [GetActualPath] on the file to get the correct name.
func (f FS) ReadFile(name string) ([]byte, error) {
	return f.fs.ReadFile(f.GetActualPath(name))
}

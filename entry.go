package hashembed

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// A wrapper around [fs.DirEntry] that keeps track of the root path.
//
// Provides mechanisms to get the root, full path, and extension of the entry.
type PathedDirEntry struct {
	entry    fs.DirEntry
	rootPath string
}

// Get the root path
func (p PathedDirEntry) RootPath() string {
	return p.rootPath
}

// Get the name of the entry
func (p PathedDirEntry) Name() string {
	return p.entry.Name()
}

// Get the name (without extension) and the extension of the entry
func (p PathedDirEntry) NameAndExtension() (string, string) {
	name := p.Name()
	ext := filepath.Ext(name)

	return strings.TrimSuffix(name, ext), strings.TrimPrefix(ext, ".")
}

// Is the entry a directory?
func (p PathedDirEntry) IsDir() bool {
	return p.entry.IsDir()
}

// Get the [fs.FileMode] of the entry
func (p PathedDirEntry) Type() fs.FileMode {
	return p.entry.Type()
}

// Get the [fs.FileInfo] of the entry
func (p PathedDirEntry) Info() (fs.FileInfo, error) {
	return p.entry.Info()
}

// Get the full path including the root to the entry
func (p PathedDirEntry) FullPath() string {
	if p.rootPath == "" {
		return p.Name()
	}
	return p.rootPath + "/" + p.Name()
}

// Create a new [PathedDirEntry] from an [fs.DirEntry]
func NewPathedDirEntry(entry fs.DirEntry, rootPath string) PathedDirEntry {
	return PathedDirEntry{
		entry,
		rootPath,
	}
}

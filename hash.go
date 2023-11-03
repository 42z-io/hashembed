package hashembed

import (
	"crypto/sha256"
	"fmt"
	"hash/crc32"
	"strconv"
)

// Interface for a FileHasher.
type FileHasher interface {
	Algorithm() string
	Hash(data []byte) (string, error)
}

// Sha256 Hasher.
type Sha256Hasher struct{}

// Return the algorithm name.
func (h Sha256Hasher) Algorithm() string {
	return "sha256"
}

// Generate a hash for the data using SHA-256.
func (h Sha256Hasher) Hash(data []byte) (string, error) {
	hash := sha256.New()
	hash.Write(data)
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// Crc32 Hasher.
type Crc32Hasher struct{}

// Return the algorithm name.
func (h Crc32Hasher) Algorithm() string {
	return "crc32"
}

// Generate a hash for the data using IEEE CRC 32.
func (h Crc32Hasher) Hash(data []byte) (string, error) {
	return strconv.FormatUint(uint64(crc32.ChecksumIEEE(data)), 16), nil
}

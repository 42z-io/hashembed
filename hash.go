package hashembed

import (
	"crypto/sha256"
	"fmt"
	"hash/crc32"
	"strconv"
)

// Interface for a FileHasher.
type FileHasher = func(data []byte) (string, error)

// Generate a hash for the data using SHA-256.
func Sha256Hasher(data []byte) (string, error) {
	hash := sha256.New()
	hash.Write(data)
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// Generate a hash for the data using IEEE CRC 32.
func Crc32Hasher(data []byte) (string, error) {
	return strconv.FormatUint(uint64(crc32.ChecksumIEEE(data)), 16), nil
}

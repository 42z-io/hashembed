package hashembed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Hasher = func([]byte) (string, error)

func TestSha256Hasher(t *testing.T) {
	hash, _ := Sha256Hasher([]byte("test"))
	assert.Equal(
		t,
		"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		hash,
		"crc32 should be calculated correctly",
	)
}

func TestCrc32Hasher(t *testing.T) {
	hash, _ := Crc32Hasher([]byte("test"))
	assert.Equal(
		t,
		"d87f7e0c",
		hash,
		"crc32 should be calculated correctly",
	)
}

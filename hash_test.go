package hashembed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha256Hasher(t *testing.T) {
	hasher := Sha256Hasher{}
	hash, _ := hasher.Hash([]byte("test"))
	assert.Equal(
		t,
		"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		hash,
		"sha256 should be calculated correctly",
	)
	assert.Equal(
		t,
		"sha256",
		hasher.Algorithm(),
		"sha256 should be the algorithm",
	)
}

func TestCrc32Hasher(t *testing.T) {
	hasher := Crc32Hasher{}
	hash, _ := hasher.Hash([]byte("test"))
	assert.Equal(
		t,
		"d87f7e0c",
		hash,
		"crc32 should be calculated correctly",
	)
	assert.Equal(
		t,
		"crc32",
		hasher.Algorithm(),
		"crc32 should be the algorithm",
	)
}

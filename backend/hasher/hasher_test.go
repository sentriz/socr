package hasher_test

import (
	"math"
	"testing"

	"github.com/matryer/is"

	"go.senan.xyz/socr/backend/hasher"
)

func TestHash(t *testing.T) {
	is := is.New(t)

	hash, err := hasher.Hash([]byte("hello"))
	is.NoErr(err)
	is.Equal(uint64(hash), uint64(2794345569481354659))

	hash, err = hasher.Hash([]byte("goodbye"))
	is.NoErr(err)
	is.Equal(uint64(hash), uint64(5515677570013980))
}

func TestFormatParse(t *testing.T) {
	is := is.New(t)

	hash := hasher.ID(5515677570013980)
	parsed, err := hasher.Parse(hash.String())
	is.NoErr(err)
	is.Equal(parsed, hash)

	hashA := uint64(math.MaxUint64)
	hash = hasher.ID(hashA)
	parsed, err = hasher.Parse(hash.String())
	is.NoErr(err)
	is.Equal(parsed, hash)

	hash = hasher.ID(0)
	parsed, err = hasher.Parse(hash.String())
	is.NoErr(err)
	is.Equal(parsed, hash)
}

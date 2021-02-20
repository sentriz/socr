package hasher

import (
	"math"
	"testing"

	"github.com/matryer/is"
)

func TestHash(t *testing.T) {
	is := is.New(t)
	hasher := &Hasher{}

	hash, err := hasher.Hash([]byte("hello"))
	is.NoErr(err)
	is.Equal(hash, uint64(2794345569481354659))

	hash, err = hasher.Hash([]byte("goodbye"))
	is.NoErr(err)
	is.Equal(hash, uint64(5515677570013980))
}

func TestFormatParse(t *testing.T) {
	is := is.New(t)
	hasher := &Hasher{}

	hash := uint64(5515677570013980)
	parsed, err := hasher.Parse(hasher.Format(hash))
	is.NoErr(err)
	is.Equal(parsed, hash)

	parsed, err = hasher.Parse(hasher.Format(math.MaxUint64))
	is.NoErr(err)
	is.Equal(parsed, uint64(math.MaxUint64))

	parsed, err = hasher.Parse(hasher.Format(0))
	is.NoErr(err)
	is.Equal(parsed, uint64(0))
}

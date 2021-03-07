package hasher

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cespare/xxhash"
)

const (
	base = 16
	bits = 64
)

type ID uint64

func (id ID) String() string {
	return strconv.FormatInt(int64(id), base)
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id ID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

func Hash(bytes []byte) (ID, error) {
	return ID(xxhash.Sum64(bytes)), nil
}

func Parse(hash string) (ID, error) {
	if hash == "" {
		return 0, fmt.Errorf("invalid hash")
	}
	i, err := strconv.ParseInt(hash, base, bits)
	return ID(i), err
}

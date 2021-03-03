package hasher

import (
	"encoding/json"
	"strconv"

	"github.com/cespare/xxhash"
)

const (
	base = 16
	bits = 64
)

type ID int64

func (id ID) String() string {
	return strconv.FormatInt(int64(id), base)
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id ID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

func Hash(bytes []byte) (int64, error) {
	return int64(xxhash.Sum64(bytes)), nil
}

func Parse(hash string) (int64, error) {
	id, err := strconv.ParseInt(hash, base, bits)
	return int64(id), err
}

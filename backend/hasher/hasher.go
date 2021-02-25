package hasher

import (
	"encoding/binary"
	"encoding/json"
	"strconv"

	"github.com/cespare/xxhash"
)

const (
	base = 16
	bits = 64
)

type ID uint64

func (id ID) String() string {
	return strconv.FormatUint(uint64(id), base)
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id ID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

type Hasher struct{}

func (h *Hasher) Hash(bytes []byte) (ID, error) {
	return ID(xxhash.Sum64(bytes)), nil
}

func (h *Hasher) Parse(hash string) (ID, error) {
	id, err := strconv.ParseUint(hash, base, bits)
	return ID(id), err
}

func (h *Hasher) Serialise(hash ID) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(hash))
	return b
}

func (h *Hasher) Deserialise(hash []byte) ID {
	return ID(binary.BigEndian.Uint64(hash))
}

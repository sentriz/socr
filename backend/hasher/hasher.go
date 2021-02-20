package hasher

import (
	"encoding/binary"
	"strconv"

	"github.com/cespare/xxhash"
)

type Hasher struct{}

func (h *Hasher) Hash(bytes []byte) (uint64, error) {
	return xxhash.Sum64(bytes), nil
}

const (
	base = 16
	bits = 64
)

func (h *Hasher) Format(hash uint64) string {
	return strconv.FormatUint(hash, base)
}

func (h *Hasher) Parse(hash string) (uint64, error) {
	return strconv.ParseUint(hash, base, bits)
}

func (h *Hasher) ToBytes(hash uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, hash)
	return b
}

func (h *Hasher) FromBytes(hash []byte) uint64 {
	return binary.BigEndian.Uint64(hash)
}

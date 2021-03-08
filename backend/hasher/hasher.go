package hasher

import (
	"strconv"

	"github.com/cespare/xxhash"
)

const (
	base = 16
	bits = 64
)

func Hash(bytes []byte) string {
	sum := xxhash.Sum64(bytes)
	format := strconv.FormatInt(int64(sum), base)
	return format
}

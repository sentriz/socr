package hasher

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cespare/xxhash"
	"github.com/jackc/pgtype"
)

const (
	base = 16
	bits = 64
)

type ID struct {
	hash uint64
}

func (id ID) String() string {
	return strconv.FormatInt(int64(id.hash), base)
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id ID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

func (id ID) Set(src interface{}) error {
	return nil
}

func (id ID) Get() interface{} {
	return nil
}

func (id ID) AssignTo(dest interface{}) error {
	return nil
}

func (id ID) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	return nil
}

func (id ID) EncodeBinary(ci *pgtype.ConnInfo, src []byte) ([]byte, error) {
	return nil, nil
}

func (id ID) DecodeText(ci *pgtype.ConnInfo, src []byte) error {
	return nil
}

func (id ID) EncodeText(ci *pgtype.ConnInfo, src []byte) ([]byte, error) {
	return nil, nil
}

func Hash(bytes []byte) (ID, error) {
	return ID{xxhash.Sum64(bytes)}, nil
}

func Parse(hash string) (ID, error) {
	if hash == "" {
		return ID{}, fmt.Errorf("invalid hash")
	}
	i, err := strconv.ParseInt(hash, base, bits)
	return ID{uint64(i)}, err
}

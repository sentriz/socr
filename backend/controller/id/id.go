package id

import (
	"math/rand"
)

const (
	idLength = 32
	idPool   = "" +
		"abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"012345679"
)

func New() string {
	bytes := make([]byte, idLength)
	for i := 0; i < len(bytes); i++ {
		bytes[i] = idPool[rand.Intn(len(idPool))]
	}

	return string(bytes)
}

package controller

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	idLength = 32
	idPool   = "" +
		"abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"012345679"
)

func IDNew() string {
	bytes := make([]byte, idLength)
	for i := 0; i < len(bytes); i++ {
		bytes[i] = idPool[rand.Intn(len(idPool))]
	}

	return string(bytes)
}

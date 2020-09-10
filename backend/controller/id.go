package controller

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	IDLength = 32
	IDPool   = "" +
		"abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"012345679"
)

func IDNew() string {
	bytes := make([]byte, IDLength)
	for i := 0; i < len(bytes); i++ {
		bytes[i] = IDPool[rand.Intn(len(IDPool))]
	}

	return string(bytes)
}

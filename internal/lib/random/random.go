package random

import (
	"math/rand"
	"time"
)

var letters = []rune(
	"abcdefghijklmnopqrstuvwxyz" + 
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789",
)

func NewRandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	runes := make([]rune, length)
	for i := range runes {
		runes[i] += letters[rnd.Intn(len(letters))]
	}

	return string(runes)
}
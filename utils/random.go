package utils

import (
	"math/rand"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandStr returns a random ASCII string of fixed length
func RandStr(n int) string {
	b := make([]rune, n)
	l := len(letters)
	for i := range b {
		b[i] = letters[rand.Intn(l)]
	}
	return string(b)
}

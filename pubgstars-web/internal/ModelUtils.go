package internal

import (
	"math/rand"
	"time"
)

func GenerateKey(n int8) string {
	var source = rand.NewSource(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	key := make([]byte, n)
	for i := range key {
		key[i] = charset[source.Int63()%int64(len(charset))]
	}
	return string(key)
}

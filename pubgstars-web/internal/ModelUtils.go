package internal

import (
	"crypto/rand"
	"math/big"
)

const keyCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func GenerateKey(n int) string {
	key := make([]byte, n)
	charsetLen := big.NewInt(int64(len(keyCharset)))
	for i := range key {
		idx, _ := rand.Int(rand.Reader, charsetLen)
		key[i] = keyCharset[idx.Int64()]
	}
	return string(key)
}

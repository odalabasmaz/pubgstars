package internal

import (
	"strings"
	"testing"
)

func TestGenerateKey_Length(t *testing.T) {
	for _, n := range []int{1, 5, 10, 32} {
		key := GenerateKey(n)
		if len(key) != n {
			t.Errorf("GenerateKey(%d): got length %d", n, len(key))
		}
	}
}

func TestGenerateKey_Charset(t *testing.T) {
	key := GenerateKey(200)
	for _, ch := range key {
		if !strings.ContainsRune(keyCharset, ch) {
			t.Errorf("GenerateKey produced unexpected character: %q", ch)
		}
	}
}

func TestGenerateKey_Uniqueness(t *testing.T) {
	a, b := GenerateKey(10), GenerateKey(10)
	if a == b {
		t.Errorf("GenerateKey produced identical keys back-to-back: %q", a)
	}
}

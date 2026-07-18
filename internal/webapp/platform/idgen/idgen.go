package idgen

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func NewHex(bytes int) (string, error) {
	buf := make([]byte, bytes)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

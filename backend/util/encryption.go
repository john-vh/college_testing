package util

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func RandString(nByte uint32) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

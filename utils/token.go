package utils

import (
	"io"
	"fmt"
	"crypto/rand"
)

func NewToken() (string, error) {
	token := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, token)
	if n != len(token) || err != nil {
		return "", err
	}
	/* variant bits */
	token[8] = token[8]&^0xc0 | 0x80
	/* pseudo-random */
	token[6] = token[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x", token), nil
}


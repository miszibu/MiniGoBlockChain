package crtpto

import (
	"crypto/sha256"
	"fmt"
)

func SHA256(str string) string{
	h := sha256.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}
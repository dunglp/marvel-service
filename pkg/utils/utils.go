package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func HashMd5(s string) (string, error) {
	h := md5.New()
	_, err := h.Write([]byte(s))
	if err != nil {
		return "", fmt.Errorf("failed to hash Md5 : %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

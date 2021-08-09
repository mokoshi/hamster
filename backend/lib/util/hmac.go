package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HmacSha256(message string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

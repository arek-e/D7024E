package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

// https://stackoverflow.com/questions/10701874/generating-the-sha-hash-of-a-string-using-golang
func Hash(data string) string {
	hasher := sha1.New()
	hasher.Write([]byte(data))
	sha := hex.EncodeToString(hasher.Sum(nil)) // Change from base64 -> hex since provided code uses that encoding type
	return sha
}

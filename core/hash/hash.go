package util

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// FileSHA256 returns the SHA256 hash of a given file
func FileSHA256(filePath string) (string, error) {
	var sha256String string

	file, err := os.Open(filePath)

	if err != nil {
		return sha256String, err
	}

	defer file.Close()

	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		return sha256String, err
	}

	hashInBytes := hash.Sum(nil)[:16]
	sha256String = hex.EncodeToString(hashInBytes)

	return sha256String, nil
}

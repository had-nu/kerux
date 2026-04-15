package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// HashFile computes the SHA-256 hash of a file efficiently
// by streaming its contents. It returns the lowercase hex digest.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil { // T2: streaming
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

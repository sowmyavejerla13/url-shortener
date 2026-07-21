package utils

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const shortCodeLength = 6

func GenerateShortCode() (string, error) {
	code := make([]byte, shortCodeLength)

	for i := range code {

		randomIndex, err := rand.Int(
			rand.Reader,
			big.NewInt(int64(len(charset))),
		)

		if err != nil {
			return "", err
		}

		code[i] = charset[randomIndex.Int64()]
	}

	return string(code), nil
}

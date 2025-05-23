package common

import (
	"crypto/rand"
	"math/big"
)

func GenerateCode() (string, error) {
	code := ""

	for i := 0; i < 5; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))

		if err != nil {
			return "", err
		}

		code += n.String()
	}

	return code, nil
}

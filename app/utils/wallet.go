package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/fatfatcocofat/rosamsoe/app/models"
)

func GenerateWalletAddress() (string, error) {
	length := 40
	byteLength := length / 2
	randomBytes := make([]byte, byteLength)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	address := hex.EncodeToString(randomBytes)

	return address, nil
}

func ValidWalletCurrency(currency string) bool {
	for _, v := range models.WalletCurrencies {
		if strings.EqualFold(v, currency) {
			return true
		}
	}

	return false
}

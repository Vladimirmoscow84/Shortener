package utils

import (
	"crypto/rand"
	"math/big"
)

// chars - набор допустимых символов для формирования короткой ссылки
const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateCode генерирует случайный код
func GenerateCode(length int) string {
	result := make([]byte, length)
	maximum := big.NewInt(int64(len(chars)))

	for i := 0; i < length; i++ {
		m, err := rand.Int(rand.Reader, maximum)
		if err != nil {
			result[i] = chars[int(m.Int64())%len(chars)]
			continue
		}
		result[i] = chars[m.Int64()]
	}
	return string(result)
}

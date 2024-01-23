package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const number = "0123456789"
const specialChar = "!@#$%^&*()-_=+[]{}|;:'<>,./?"

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int) int {
	return min + r.Intn(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomName() string {
	return RandomString(RandomInt(3, 6))
}

func RandomEmail() string {
	return RandomName() + "@gmail.com"
}

func RandomPassword() string {
	passwordLen := RandomInt(6, 12)
	var sb strings.Builder
	char := alphabet + number + specialChar
	k := len(char)

	for i := 0; i < passwordLen; i++ {
		c := char[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

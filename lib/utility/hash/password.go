package hash

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func IsValidPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func verifyPassword(s string) bool {
	letters := 0
	var number, upper, special, sevenOrMore bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			return false
		}
	}
	sevenOrMore = letters >= 7
	return number && upper && special && sevenOrMore
}

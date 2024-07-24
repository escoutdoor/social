package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPw(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePw(pw string, hashedPw string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(pw)); err != nil {
		return false
	}
	return true
}

package utils

import (
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

func BcryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func BcryptCompare(password, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

func VerifyPasswordComplexity(password string) bool {
	var (
		hasMinLength = false
		hasUpper     = false
		hasLower     = false
		hasNumber    = false
		hasSpecial   = false
	)

	if len(password) >= 8 {
		hasMinLength = true
	}
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	return hasMinLength && hasUpper && hasLower && hasNumber && hasSpecial
}

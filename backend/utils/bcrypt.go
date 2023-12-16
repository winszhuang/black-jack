package utils

import "golang.org/x/crypto/bcrypt"

const cost = 8

func GenerateFromPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func CompareHashAndPassword(hashPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashPassword, password)
}

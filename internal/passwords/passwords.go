package passwords

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashFromPassword(password string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPasswordBytes), err
}

func ComparePasswordWithHash(passwordAttempt, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordAttempt))
}

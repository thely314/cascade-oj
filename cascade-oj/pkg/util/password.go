package util

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// use bcrypt to generate hashed password
func GenerateHashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password is empty")
	}
	// use bcrypt default cost
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// use CompareHashAndPassword to verify password
// params:
//   - password: the MD5 hashed password to verify
//   - stored: the stored hashed password (bcrypt hashed MD5 password)
//
// returns:
//   - bool: true if password matches, false otherwise
func VerifyPassword(password, stored string) bool {
	if password == "" || stored == "" {
		return false
	}
	// use bcrypt secure comparison: returns nil on success
	err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(password))
	return err == nil
}

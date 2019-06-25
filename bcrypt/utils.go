package bcrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Utilities implements the HashUtilities interface to provide hashing operations
type Utilities struct {
	HashCost int
}

// HashAndSalt computes a hash given an input string and a hash cost
func (u Utilities) HashAndSalt(s string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), u.HashCost)
	if err != nil {
		return "", errors.New("Failed to hash string: " + err.Error())
	}
	return string(bytes), nil
}

// CompareHashWithString checks if a given hash is equivalent to a string should
// the string be hashed
func (u Utilities) CompareHashWithString(hash string, s string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(s)) == nil
}

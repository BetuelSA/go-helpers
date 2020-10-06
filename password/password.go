package password

import (
	"golang.org/x/crypto/bcrypt"
)

type password struct{}

// PasswordSvc --
type Service interface {
	Hash(password string) (string, error)
	CheckPassword(password, hash string) error
}

// NewPasswordSvc creates a new fake password
func NewService() Service {
	return &password{}
}

// Hash returns the hash for a password
func (p *password) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash), err
}

// CheckPassword compares a password and its hash, validating it
func (p *password) CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

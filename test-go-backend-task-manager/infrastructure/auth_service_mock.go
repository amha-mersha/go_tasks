package infrastructure

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Define an interface for the functions you want to mock
type AuthService interface {
	HashPassword(password string) (string, error)
	ValidatePassword(hashedPassword, password string) error
	CreateJWTToken(username string, role string, timeDuration time.Duration) (string, error)
	ParseJWTToken(token string) (*jwt.Token, error)
}

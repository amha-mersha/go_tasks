package infrastructure

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaim struct {
	Username string `json:"username"`
	UserRole string `json:"userrole"`
	jwt.RegisteredClaims
}

func CreateJWTToken(username string, userrole string, timeDuration time.Duration) (string, error) {
	claim := UserClaim{
		Username: username,
		UserRole: userrole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(timeDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	jwtToken, err := token.SignedString([]byte(os.Getenv("SIGNITURE_SECRET")))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

func ParseJWTToken(token string) (*jwt.Token, error) {
	retrivedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("SIGNITURE_SECRET")), nil
	})

	if err != nil {
		return &jwt.Token{}, err
	}
	if !retrivedToken.Valid {
		return &jwt.Token{}, fmt.Errorf("Invalid token")
	}
	return retrivedToken, nil
}

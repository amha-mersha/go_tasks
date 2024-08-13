package infrastructure

import "github.com/golang-jwt/jwt/v5"

type UserClaim struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserRole string `json:"userrole"`
	jwt.RegisteredClaims
}

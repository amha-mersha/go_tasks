package infrastructure

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleWare(validRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Authorization header is required"})
			ctx.Abort()
			return
		}
		authTokens := strings.Split(authHeader, " ")
		if len(authTokens) != 2 || authTokens[0] != "Bearer" {
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid authorization header"})
			ctx.Abort()
			return
		}
		token, err := ParseJWTToken(authTokens[1])
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid token"})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid token claims"})
			ctx.Abort()
			return
		}
		expirationDate, ok := claims["exp"].(float64)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid token, Token expiration date not found"})
			ctx.Abort()
			return
		}
		issuedDate, ok := claims["iat"].(float64)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid token, Token issued date not found"})
			ctx.Abort()
			return
		}

		if time.Now().After(time.Unix(int64(expirationDate), 0)) || time.Now().Before(time.Unix(int64(issuedDate), 0)) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "Token expired"})
			ctx.Abort()
			return
		}

		retrivedRole, ok := token.Claims.(jwt.MapClaims)["role"]
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid token, Role of the user is not found"})
			ctx.Abort()
			return
		}
		valid := false
		for _, role := range validRoles {
			if role == retrivedRole {
				valid = true
				break
			}
		}
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized access"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

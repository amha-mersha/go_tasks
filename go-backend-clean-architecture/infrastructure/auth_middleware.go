package infrastructure

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Authorization header is required"})
		}
		authTokens := strings.Split(authHeader, " ")
		if len(authTokens) != 2 || authTokens[0] != "Bearer" {
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid authorization header"})
		}

		token, ok := jwt.ParseWithClaims(authTokens[1])

		ctx.Next()
	}
}

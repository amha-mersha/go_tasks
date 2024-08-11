package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/amha-mersha/go_tasks/authentication_authorization_go_mongo/data"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// getting the Authorization header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Authorization authHeader is required"})
			c.Abort()
			return
		}
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Invalid authorization header"})
			c.Abort()
			return
		}

		// checking for the correctness of the token
		claim := &data.UserCustomClaim{}
		token, err := jwt.ParseWithClaims(authParts[1], claim, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method :%v", t.Header["alg"])
			}
			return []byte(os.Getenv("SIGNITURE_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": data.InternalServerError + " " + err.Error()})
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid token"})
			c.Abort()
			return
		}

		// cheking if the time range of the claim haven't expired
		if time.Now().After(claim.ExpiresAt.Time) || time.Now().Before(claim.IssuedAt.Time) {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Token expired"})
			c.Abort()
			return
		}

		// cheking if the username exists in the database
		// check role and the route they are trying to access
		user, err := data.GetUserByUsername(claim.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": data.UserNotFound})
			c.Abort()
			return
		}
		if claim.Username != user.Username {
			c.JSON(http.StatusBadRequest, gin.H{"Error": data.UserNotFound})
			c.Abort()
			return
		}

		if !PathMap[c.Request.Method][claim.UserRole] {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "you are not authorized to use this route"})
			c.Abort()
			return
		}

		c.Next()
	}

}

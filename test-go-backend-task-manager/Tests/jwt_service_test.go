package tests

import (
	"os"
	"testing"
	"time"

	"github.com/amha-mersha/go_tasks/test-go-backend-task-manager/infrastructure"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type JWTTestSuite struct {
	suite.Suite
	signatureSecret string
}

func (suite *JWTTestSuite) SetupSuite() {
	// Set up environment variable for testing
	suite.signatureSecret = "test_secret"
	os.Setenv("SIGNITURE_SECRET", suite.signatureSecret)
}

func (suite *JWTTestSuite) TearDownSuite() {
	// Clean up environment variable
	os.Unsetenv("SIGNITURE_SECRET")
}

func (suite *JWTTestSuite) TestCreateJWTToken_Success() {
	username := "testuser"
	role := "admin"
	duration := time.Minute * 10

	token, err := infrastructure.CreateJWTToken(username, role, duration)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)

	// Optional: Decode the token to verify the claims
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(suite.signatureSecret), nil
	})
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), parsedToken.Valid)
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), username, claims["username"])
	assert.Equal(suite.T(), role, claims["role"])
}

func (suite *JWTTestSuite) TestParseJWTToken_Success() {
	username := "testuser"
	role := "admin"
	duration := time.Minute * 10

	token, err := infrastructure.CreateJWTToken(username, role, duration)
	assert.NoError(suite.T(), err)

	parsedToken, err := infrastructure.ParseJWTToken(token)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), username, claims["username"])
	assert.Equal(suite.T(), role, claims["role"])
}

func (suite *JWTTestSuite) TestParseJWTToken_InvalidToken() {
	invalidToken := "invalid.token.here"
	_, err := infrastructure.ParseJWTToken(invalidToken)
	assert.Error(suite.T(), err)
}

func (suite *JWTTestSuite) TestParseJWTToken_ExpiredToken() {
	username := "testuser"
	role := "admin"
	expiredDuration := -time.Minute * 10

	token, err := infrastructure.CreateJWTToken(username, role, expiredDuration)
	assert.NoError(suite.T(), err)

	parsedToken, err := infrastructure.ParseJWTToken(token)
	assert.Error(suite.T(), err)
	assert.False(suite.T(), parsedToken.Valid)
}

func TestJWTTestSuite(t *testing.T) {
	suite.Run(t, new(JWTTestSuite))
}

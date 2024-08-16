package tests

import (
	"testing"

	"github.com/amha-mersha/go_tasks/test-go-backend-task-manager/infrastructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PasswordTestSuite struct {
	suite.Suite
}

func (suite *PasswordTestSuite) TestHashPassword_Success() {
	password := "testpassword"
	hashedPassword, err := infrastructure.HashPassword(password)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hashedPassword)
}

func (suite *PasswordTestSuite) TestValidatePassword_Success() {
	password := "testpassword"
	hashedPassword, err := infrastructure.HashPassword(password)
	assert.NoError(suite.T(), err)

	err = infrastructure.ValidatePassword(hashedPassword, password)
	assert.NoError(suite.T(), err)
}

func (suite *PasswordTestSuite) TestValidatePassword_Failure() {
	password := "testpassword"
	wrongPassword := "wrongpassword"
	hashedPassword, err := infrastructure.HashPassword(password)
	assert.NoError(suite.T(), err)

	err = infrastructure.ValidatePassword(hashedPassword, wrongPassword)
	assert.Error(suite.T(), err)
}

func TestPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordTestSuite))
}

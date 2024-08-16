package infrastructure

import "golang.org/x/crypto/bcrypt"

// type auth_Service struct {
// }
//
// func NewAuthService() auth_Service {
// 	return auth_Service{}
// }

func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashPassword), err
}

func ValidatePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

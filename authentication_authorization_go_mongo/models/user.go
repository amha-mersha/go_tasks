package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type IssuedUser struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

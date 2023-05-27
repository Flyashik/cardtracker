package models

type User struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

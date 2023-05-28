package models

type User struct {
	Id       string `json:"user_id"`
	Name     string `json:"name"`
	Code     int    `json:"code"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

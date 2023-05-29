package models

type UserPhone struct {
	User   User  `json:"user"`
	Phones []int `json:"phones"`
}

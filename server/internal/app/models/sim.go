package models

type SimInfo struct {
	Id          int
	PhoneNumber string `json:"phone_number"`
	Operator    string `json:"operator"`
}

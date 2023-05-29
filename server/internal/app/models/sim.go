package models

type SimInfo struct {
	Id          int    `json:"sim_id"`
	PhoneId     *int   `json:"phone_id"`
	PhoneNumber string `json:"phone_number"`
	Operator    string `json:"operator"`
}

package models

type SimInfo struct {
	Imei        string  `json:"imei"`
	PhoneNumber *string `json:"phone_number"`
	Operator    *string `json:"operator"`
}

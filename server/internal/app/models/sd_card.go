package models

type SdInfo struct {
	Id               int    `json:"sd_card_id"`
	PhoneId          *int   `json:"phone_id"`
	SdManufacturerId string `json:"sd_manufacturer_id"`
	SerialNo         string `json:"serial_no"`
	TotalSpace       int    `json:"total_space"`
	UsedSpace        int    `json:"used_space"`
	FreeSpace        int    `json:"free_space"`
}

package models

type SdInfo struct {
	Id               int
	SdManufacturerId string `json:"sd_manufacturer_id"`
	SerialNo         string `json:"serial_no"`
	TotalSpace       int    `json:"total_space"`
	UsedSpace        int    `json:"used_space"`
	FreeSpace        int    `json:"free_space"`
}

package models

type SdInfo struct {
	SdManufacturer string `json:"sd_manufacturer"`
	TotalSpace     int    `json:"total_space"`
	UsedSpace      int    `json:"used_space"`
	FreeSpace      int    `json:"free_space"`
}

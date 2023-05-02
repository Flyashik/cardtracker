package models

type Info struct {
	Phone   Phone     `json:"phone_info"`
	SimInfo []SimInfo `json:"sim_info"`
	SdInfo  []SdInfo  `json:"sd_info"`
	AuthID  uint      `json:"authorization_id"`
}

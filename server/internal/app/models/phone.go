package models

type Phone struct {
	Id             int
	Manufacturer   string   `json:"manufacturer"`
	ModelTag       string   `json:"model_tag"`
	ModelNumber    string   `json:"model_number"`
	OsVersion      string   `json:"os_version"`
	ApiVersion     string   `json:"api_version"`
	Cpu            string   `json:"cpu"`
	Firmware       string   `json:"firmware"`
	Bootloader     string   `json:"bootloader"`
	SupportedArchs []string `json:"supported_archs"`
	SimSlots       int
	SdSlots        int
}

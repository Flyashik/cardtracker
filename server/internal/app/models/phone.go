package models

type Phone struct {
	Manufacturer   string   `json:"manufacturer"`
	Model          string   `json:"model"`
	ModelNumber    string   `json:"model_number"`
	OsVersion      string   `json:"os_version"`
	ApiVersion     string   `json:"api_version"`
	Cpu            string   `json:"cpu"`
	Firmware       string   `json:"firmware"`
	SupportedArchs []string `json:"supported_archs"`
}

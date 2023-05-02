package storage

type DbConfig struct {
	DbURL string `toml:"db_url"`
}

func NewConfig() *DbConfig {
	return &DbConfig{}
}

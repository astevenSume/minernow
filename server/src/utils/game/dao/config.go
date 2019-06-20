package dao

type Config struct {
	OfflineInterval int64 // game user offline interval
}

func NewConfig() *Config {
	return &Config{}
}

var GameConfig = NewConfig()

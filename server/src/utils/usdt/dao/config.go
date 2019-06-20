package dao

type Config struct {
	OnchainDataRegionNum uint64 //how many regions on chain data be divided to.
}

func NewConfig() *Config {
	return &Config{}
}

var DaoConfig = NewConfig()

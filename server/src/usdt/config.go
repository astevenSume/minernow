package usdt

type RpcConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Config struct {
	RpcConfig
	Precision         int
	OfflineTxPort     int
	MinFee, MaxFee    int //fee value should be in [MinFee, MaxFee]
	UnitFee			  int //default uint fee
	Symbol            string
	ConfirmationLimit int
	FeeMode           int
	SyncFrequency     int64 // the frequency limit for user sync transactions from chain.
	PlatformUaid	  uint64
}

var UsdtConfig Config

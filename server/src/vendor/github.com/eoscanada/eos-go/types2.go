package eos

type SignedTransactionResponse struct {
	Expiration            string        `json:"expiration"`
	RefBlockNum           int           `json:"ref_block_num"`
	RefBlockPrefix        int           `json:"ref_block_prefix"`
	MaxNetUsageWords      int           `json:"max_net_usage_words"`
	MaxCpuUsageMs         int           `json:"max_cpu_usage_ms"`
	DelaySec              int           `json:"delay_sec"`
	ContextFreeActions    []interface{} `json:"context_free_actions"`
	Actions               []*Action     `json:"actions"`
	TransactionExtensions []interface{} `json:"transaction_extensions"`
	Signatures            []string      `json:"signatures"`
	ContextFreeData       []interface{} `json:"context_free_data"`
}

type TransactionRespV16 struct {
	Id  string `json:"id"`
	Trx struct {
		Receipt struct {
			Status        string        `json:"status"`
			CpuUsageUs    int           `json:"cpu_usage_us"`
			NetUsageWords int           `json:"net_usage_words"`
			Trx           []interface{} `json:"trx"`
		} `json:"receipt"`
		Trx TransactionTrxTrxV16 `json:"trx"`
	} `json:"trx"`
	BlockTime             string `json:"block_time"`
	BlockNum              uint32 `json:"block_num"`
	LastIrreversibleBlock uint32 `json:"block_id"`
	//Traces      []string               `json:"traces"`
}

type TransactionTrxTrxV16 struct {
	Expiration            string                       `json:"expiration"`
	RefBlockNum           int                          `json:"ref_block_num"`
	RefBlockPrefix        int                          `json:"ref_block_prefix"`
	MaxNetUsageWords      int                          `json:"max_net_usage_words"`
	MaxCpuUsageMs         int                          `json:"max_cpu_usage_ms"`
	DelaySec              int                          `json:"delay_sec"`
	ContextFreeActions    []interface{}                `json:"context_free_actions"`
	Actions               []TransactionTrxTrxActionV16 `json:"actions"`
	TransactionExtensions []interface{}                `json:"transaction_extensions"`
	Signatures            []string                     `json:"signatures"`
	ContextFreeData       []interface{}                `json:"context_free_data"`
}
type TransactionTrxTrxActionV16 struct {
	Account       string                         `json:"account"`
	Name          string                         `json:"name"`
	Authorization []interface{}                  `json:"authorization"`
	Data          TransactionTrxTrxActionDataV16 `json:"data"`
	HexData       string                         `json:"hex_data"`
}
type TransactionTrxTrxActionDataV16 struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Quantity string `json:"quantity"`
	Memo     string `json:"memo"`
}

//push transaction
type PushTransaction struct {
	Compression string                 `json:"compression"`
	Transaction map[string]interface{} `json:"transaction"`
	Signatures  []string               `json:"signatures"`
}

// push transaction response
type PushTransactionResponse struct {
	TransactionId string                           `json:"transaction_id"`
	Processed     PushTransactionResponseProcessed `json:"processed"`
}

type PushTransactionResponseProcessed struct {
	Id              string `json:"id"`
	BlockNum        uint32 `json:"block_num"`
	BlockTime       string `json:"block_time"`
	ProducerBlockId string `json:"producer_block_id"`
	Receipt         struct {
		Status           TransactionStatus `json:"status"`
		CPUUsageMicrosec int               `json:"cpu_usage_us"`
		NetUsageWords    int               `json:"net_usage_words"`
	} `json:"receipt"`
	Elapsed      int           `json:"elapsed"`
	NetUsage     int           `json:"net_usage"`
	Scheduled    bool          `json:"scheduled"`
	ActionTraces []interface{} `json:"action_traces"`
}

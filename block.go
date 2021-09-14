package bitcoindata

type Block struct {
	Hash         string `json:"hash"`
	Version      int    `json:"ver"`
	PrevBlock    string `json:"prev_block"`
	MerkleRoot   string `json:"mrkl_root"`
	Time         int    `json:"time"`
	Bits         int    `json:"bits"`
	Nonce        int    `json:"nonce"`
	TXCount      int    `json:"n_tx"`
	Size         int    `json:"size"`
	BlockIndex   int    `json:"block_index"`
	MainChain    bool   `json:"main_chain"`
	Height       int    `json:"height"`
	ReceivedTime int    `json:"received_time"`
	RelayedBy    string `json:"relayed_by"`
	// I will omite transactions to implement for the future
}

type APIResponse struct {
	Blocks []Block `json:"blocks"`
}

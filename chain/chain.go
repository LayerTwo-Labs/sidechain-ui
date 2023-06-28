package chain

type ChainData struct {
	Regtest int    `json:"regtest"`
	Port    int    `json:"rpcport"`
	RPCUser string `json:"rpcuser"`
	RPCPass string `json:"rpcpassword"`
	Slot    int    `json:"slot,omitempty"`
}

package models

type Wallet struct {
	Address string `json:"address"`
	Balance int64  `json:"balance"`
	Chain   string `json:"chain"`
	Token   string `json:"token"`
}

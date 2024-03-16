package models

type WalletMessage struct {
	Address string `json:"address"`
	Balance int64  `json:"balance"`
	Chain   string `json:"chain"`
	Token   string `json:"token"`
}

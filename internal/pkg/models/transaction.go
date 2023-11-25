package models

type Transaction struct {
	Hash     string  `json:"hash"`
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
	Chain    string  `json:"chain"`
	Token    string  `json:"token"`
}

package models

type Transaction struct {
	Txid     string             `json:"txid"`
	InputTxn []InputTransaction `json:"vin"`
}

type InputTransaction struct {
	Txid string `json:"txid"`
}

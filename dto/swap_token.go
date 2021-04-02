package dto

type SwapToken struct {
	FromWalletId string  `json:"fromWalletId"`
	ToWalletId   string  `json:"toWalletId"`
	Amount       float64 `json:"amount"`
}

package model

type Merchant struct {
	MerchantID      int `json:"merchant_id"`
	MerchantUserID  int `json:"user_id"`
	MerchantBalance int `json:"balance"`
}

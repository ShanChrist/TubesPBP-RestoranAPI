package model

type Promo struct {
	PromoKode       string  `json:"kode_promo"`
	PromoDeskripsi  string  `json:"deskripsi_promo"`
	PromoPersentase float32 `json:"persentase"`
}

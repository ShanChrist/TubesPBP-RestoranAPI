package model

type Restoran struct {
	RestoranID         int    `json:"restoran_id`
	RestoranNama       string `json:"nama_restoran`
	RestoranJamBuka    string `json:"jam_buka`
	RestoranJamTutup   string `json:"jam_tutup`
	RestoranAlamat     string `json:"alamat_restoran`
	RestoranHariBuka   string `json:"hari_buka`
	RestoranMerchantID int    `json:"merchant_id"`
}

type DetailRestoran struct {
	RestoranID       int       `json:"restoran_id`
	RestoranNama     string    `json:"nama_restoran`
	RestoranJamBuka  string    `json:"jam_buka`
	RestoranJamTutup string    `json:"jam_tutup`
	RestoranAlamat   string    `json:"alamat_restoran`
	RestoranHariBuka string    `json:"hari_buka`
	MakananData      []Makanan `json:"data"`
}

package model

type DetailTransaksi struct {
	DetailTransaksiID          int `json:"detail_transaksi_id"`
	DetailTransaksiTransaksiID int `json:"transaksi_id"`
	DetailTransaksiMakananID   int `json:"makanan_id"`
	DetailTransaksiKuantitas   int `json:"kuantitas"`
}

type DetailOrderTransaksi struct {
	DetailTransaksiMakananID int `json:"makanan_id"`
	DetailTransaksiKuantitas int `json:"kuantitas"`
}

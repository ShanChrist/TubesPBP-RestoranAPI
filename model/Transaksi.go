package model

type Transaksi struct {
	TransakasiID         int             `json:"transaksi_id"`
	TransakasiRestoranID int             `json:"restoran_id"`
	TransakasiUserID     int             `json:"user_id"`
	TransakasiHargaTotal int             `json:"harga_total"`
	TransakasiStatus     TransaksiStatus `json:"status"`
}

type DetailMemberTransaksi struct {
	TransakasiID         int                   `json:"transaksi_id"`
	TransakasiHargaTotal int                   `json:"harga_total"`
	TransakasiStatus     DetailTransaksiStatus `json:"status"`
}

type TransaksiStatus struct {
	TransaksiStatusID   int    `json:"id"`
	TransaksiStatusName string `json:"name"`
}

type DetailTransaksiStatus struct {
	TransaksiStatusName string `json:"name"`
}

type DetailTransaksiOrder struct {
	TransakasiID         int                  `json:"transaksi_id"`
	TransakasiHargaTotal int                  `json:"harga_total"`
	TransaksiDetailOrder DetailOrderTransaksi `json:"data"`
}

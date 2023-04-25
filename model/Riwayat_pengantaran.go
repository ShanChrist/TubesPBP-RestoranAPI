package model

type RiwayatPengantaran struct {
	RiwayatPengantaranID           int `json:"pengantaran_id"`
	RiwayatPengantaranCourrierID   int `json:"courrier_id"`
	RiwayatPengantaranTransakasiID int `json:"transaksi_id"`
}

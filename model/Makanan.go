package model

type Makanan struct {
	MakananID         int    `json:"makanan_id"`
	MakananRestoranID int    `json:"restoran_id"`
	MakananNama       string `json:"nama_makanan"`
	MakananHarga      int    `json:"harga_makanan"`
}

type DetailMakanan struct {
	MakananID    int    `json:"makanan_id"`
	MakananNama  string `json:"nama_makanan"`
	MakananHarga int    `json:"harga_makanan"`
}

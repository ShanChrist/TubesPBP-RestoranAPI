package model

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type LihatRestoranResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []Restoran `json:"data"`
}

type LihatSpecificRestoranResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    []DetailRestoran `json:"data"`
}

type LihatStatusOrderResponse struct {
	Status  int                     `json:"status"`
	Message string                  `json:"message"`
	Data    []DetailMemberTransaksi `json:"data"`
}

type LihatPromoResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    []Promo `json:"data"`
}

type LihatMakananResponse struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    []DetailMakanan `json:"data"`
}

type LihatSpecificRestoran struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    DetailRestoran `json:"data"`
}

type LihatDetailOrderResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    []DetailTransaksiOrder `json:"data"`
}

type LihatRiwayatPengantaranResponse struct {
	Status  int                  `json:"status"`
	Message string               `json:"message"`
	Data    []RiwayatPengantaran `json:"data"`
}

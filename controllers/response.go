package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Tubes/model"
)

func sendResponse(w http.ResponseWriter, status int, message string) {
	var response model.Response
	response.Status = status
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func lihatRestoranResponse(w http.ResponseWriter, restorans []model.Restoran) {
	var response model.LihatRestoranResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = restorans
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func lihatSpecificRestoranResponse(w http.ResponseWriter, restorans []model.DetailRestoran) {
	var response model.LihatSpecificRestoranResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = restorans
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LihatStatusOrderResponse(w http.ResponseWriter, transaksis []model.DetailMemberTransaksi) {
	var response model.LihatStatusOrderResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = transaksis
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func lihatPromoResponse(w http.ResponseWriter, promos []model.Promo) {
	var response model.LihatPromoResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = promos
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func lihatMakananResponse(w http.ResponseWriter, makanans []model.DetailMakanan) {
	var response model.LihatMakananResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = makanans
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LihatDetailOrderResponse(w http.ResponseWriter, transaksis []model.DetailTransaksiOrder) {
	var response model.LihatDetailOrderResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = transaksis
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LihatRiwayatPengantaranResponse(w http.ResponseWriter, pengantarans []model.RiwayatPengantaran) {
	var response model.LihatRiwayatPengantaranResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = pengantarans
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func lihatProfileResponse(w http.ResponseWriter, user model.LihatProfile) {
	var response model.LihatProfileResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Tubes/model"
	"github.com/gorilla/mux"
)

func InsertPromo(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	promo := model.Promo{
		PromoKode:      r.Form.Get("kode_promo"),
		PromoDeskripsi: r.Form.Get("deskripsi_promo"),
	}

	persentasePromoStr := r.Form.Get("persentase")
	persentasePromo, err := strconv.ParseFloat(persentasePromoStr, 32)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	promo.PromoPersentase = float32(persentasePromo)

	if promo.PromoKode == "" || promo.PromoDeskripsi == "" || promo.PromoPersentase == 0 {
		http.Error(w, "False || Invalid input", http.StatusBadRequest)
		return
	}

	_, userID, _, _ := validateTokenFromCookies(r)

	var userType int
	row := db.QueryRow("SELECT user_type FROM users WHERE user_id=?", userID)
	if err := row.Scan(&userType); err != nil {
		log.Println("Error scanning user count:", err)
		return
	}

	if userID <= 0 || userType != 0 {
		sendResponse(w, 400, "Access Denied || Not Admin")
		return
	}

	// Insert promo into database
	_, err = db.Exec("INSERT INTO promo (kode_promo, deskripsi_promo, persentase) VALUES (?, ?, ?)", promo.PromoKode, promo.PromoDeskripsi, promo.PromoPersentase)

	if err != nil {
		log.Println("Error inserting promo:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sendResponse(w, 200, "Success Added Promo")
}

func DeletePromo(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	params := mux.Vars(r)
	kode_promo := params["kode_promo"]

	_, userID, _, _ := validateTokenFromCookies(r)

	var userType int
	row := db.QueryRow("SELECT user_type FROM users WHERE user_id=?", userID)
	if err := row.Scan(&userType); err != nil {
		log.Println("Error scanning user count:", err)
		return
	}

	if userID <= 0 || userType != 0 {
		sendResponse(w, 400, "Access Denied || Not Admin")
		return
	}

	// Delete menu item from database
	result, err := db.Exec("DELETE FROM promo WHERE kode_promo=?", kode_promo)

	if err != nil {
		log.Println("Error deleting makanan:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Println("Error getting rows affected:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		sendResponse(w, 404, "Promo not found")
		return
	}
	sendResponse(w, 200, "Promo deleted successfully")
}

func Lihat(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	_, userID, _, _ := validateTokenFromCookies(r)

	var userType int
	row := db.QueryRow("SELECT user_type FROM users WHERE user_id=?", userID)
	if err := row.Scan(&userType); err != nil {
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Login as Member First")
		return
	}

	if userID <= 0 || userType != 1 {
		sendResponse(w, 400, "Access Denied || Login as Member First")
		return
	}

	query := "select * from restoran"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var restoran model.Restoran
	var restorans []model.Restoran
	for rows.Next() {
		if err := rows.Scan(&restoran.RestoranID, &restoran.RestoranNama, &restoran.RestoranJamBuka, &restoran.RestoranJamTutup, &restoran.RestoranAlamat, &restoran.RestoranHariBuka, &restoran.RestoranMerchantID); err != nil {
			log.Println(err)
			return
		} else {
			restorans = append(restorans, restoran)
		}
	}

	if len(restorans) < 1 {
		sendResponse(w, 400, "Error Array Size Not Correct")
	} else {
		lihatRestoranResponse(w, restorans)
	}
}

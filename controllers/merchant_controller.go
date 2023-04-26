package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Tubes/model"
	"github.com/gorilla/mux"
)

func RestoranRegister(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	restoran := model.Restoran{
		RestoranNama:     r.Form.Get("nama_restoran"),
		RestoranJamBuka:  r.Form.Get("jam_buka"),
		RestoranJamTutup: r.Form.Get("jam_tutup"),
		RestoranAlamat:   r.Form.Get("alamat_restoran"),
		RestoranHariBuka: r.Form.Get("hari_buka"),
	}

	if restoran.RestoranNama == "" || restoran.RestoranJamBuka == "" || restoran.RestoranJamTutup == "" || restoran.RestoranAlamat == "" || restoran.RestoranHariBuka == "" {
		http.Error(w, "False || Invalid input", http.StatusBadRequest)
		return
	}

	_, userID, _, _ := validateTokenFromCookies(r)

	var userType int
	row := db.QueryRow("SELECT user_type FROM users WHERE user_id=?", userID)
	if err := row.Scan(&userType); err != nil {
		sendResponse(w, 400, "Login as Merchant First")
		return
	}

	if userID <= 0 || userType != 2 {
		sendResponse(w, 400, "Login as Merchant to Create Restoran")
		return
	}

	var merchantID int
	row = db.QueryRow("SELECT merchant_id FROM merchant WHERE user_id=?", userID)

	if err := row.Scan(&merchantID); err != nil {
		log.Println("Error scanning user count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var count int
	// Check if restoran merchant already make one
	row = db.QueryRow("SELECT COUNT(*) FROM restoran WHERE merchant_id=?", merchantID)

	if err := row.Scan(&count); err != nil {
		log.Println("Error scanning user count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		sendResponse(w, 400, "You already created restoran")
		return
	}

	// Insert restaurant into database
	_, err = db.Exec("INSERT INTO restoran (nama_restoran, jam_buka, jam_tutup, alamat_restoran, hari_buka, merchant_id) VALUES (?, ?, ?, ?, ?, ?)", restoran.RestoranNama, restoran.RestoranJamBuka, restoran.RestoranJamTutup, restoran.RestoranAlamat, restoran.RestoranHariBuka, merchantID)

	if err != nil {
		log.Println("Error inserting restaurant:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sendResponse(w, 200, "Success Added Restoran")
}

func InsertMakanan(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	makanan := model.Makanan{
		MakananNama: r.Form.Get("nama_makanan"),
	}

	if makanan.MakananNama == "" {
		http.Error(w, "False || Invalid input", http.StatusBadRequest)
		return
	}

	harga_makananStr := r.Form.Get("harga_makanan")
	harga_makanan, err := strconv.Atoi(harga_makananStr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	makanan.MakananHarga = harga_makanan

	_, userID, _, _ := validateTokenFromCookies(r)

	var userType int
	row := db.QueryRow("SELECT user_type FROM users WHERE user_id=?", userID)
	if err := row.Scan(&userType); err != nil {
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Login as Merchant First")
		return
	}

	if userID <= 0 || userType != 2 {
		sendResponse(w, 400, "Access Denied || Login as Merchant First")
		return
	}

	var merchantID int
	row = db.QueryRow("SELECT merchant_id FROM merchant WHERE user_id=?", userID)

	if err := row.Scan(&merchantID); err != nil {
		log.Println("Error scanning user count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var restoranID int
	row = db.QueryRow("SELECT restoran_id FROM restoran WHERE merchant_id=?", merchantID)

	if err := row.Scan(&restoranID); err != nil {
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Access Denied || Login as Merchant or Create Restoran First")
		return
	}

	if restoranID < 0 {
		sendResponse(w, 400, "Access Denied || Login as Merchant or Create Restoran First")
		return
	}

	// Insert restaurant into database
	_, err = db.Exec("INSERT INTO makanan (restoran_id, nama_makanan, harga_makanan) VALUES (?, ?, ?)", restoranID, makanan.MakananNama, makanan.MakananHarga)

	if err != nil {
		log.Println("Error inserting makanan:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sendResponse(w, 200, "Success Added Menu")
}

func DeleteMakanan(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	params := mux.Vars(r)
	makananID := params["makanan_id"]

	_, userID, _, _ := validateTokenFromCookies(r)

	var userType int
	row := db.QueryRow("SELECT user_type FROM users WHERE user_id=?", userID)
	if err := row.Scan(&userType); err != nil {
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Login as Merchant First")
		return
	}

	if userID <= 0 || userType != 2 {
		sendResponse(w, 400, "Access Denied || Login as Merchant First")
		return
	}

	var merchantID int
	row = db.QueryRow("SELECT merchant_id FROM merchant WHERE user_id=?", userID)

	if err := row.Scan(&merchantID); err != nil {
		log.Println("Error scanning user count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var restoranID int
	row = db.QueryRow("SELECT restoran_id FROM restoran WHERE merchant_id=?", merchantID)

	if err := row.Scan(&restoranID); err != nil {
		log.Println("Error scanning user count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if restoranID < 0 {
		sendResponse(w, 400, "Access Denied || Login as Merchant or Create Restoran First")
		return
	}

	// Delete menu item from database
	result, err := db.Exec("DELETE FROM makanan WHERE makanan_id=? AND restoran_id=?", makananID, restoranID)

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
		sendResponse(w, 404, "Menu item not found")
		return
	}

	sendResponse(w, 200, "Menu item deleted successfully")
}

func LihatMakanan(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	_, userID, _, _ := validateTokenFromCookies(r)

	var userType int
	row := db.QueryRow("SELECT user_type FROM users WHERE user_id=?", userID)
	if err := row.Scan(&userType); err != nil {
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Login as Merchant First")
		return
	}

	if userID <= 0 || userType != 2 {
		sendResponse(w, 400, "Access Denied || Login as Merchant First")
		return
	}

	var merchantID int
	row = db.QueryRow("SELECT merchant_id FROM merchant WHERE user_id=?", userID)

	if err := row.Scan(&merchantID); err != nil {
		log.Println("Error scanning user count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var restoranID int
	row = db.QueryRow("SELECT restoran_id FROM restoran WHERE merchant_id=?", merchantID)

	if err := row.Scan(&restoranID); err != nil {
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Access Denied || Login as Merchant or Create Restoran First")
		return
	}

	query := `SELECT m.makanan_id, m.nama_makanan, m.harga_makanan FROM makanan m JOIN restoran r ON m.restoran_id = r.restoran_id JOIN merchant me ON r.merchant_id = me.merchant_id WHERE me.merchant_id = ? AND r.restoran_id = ?`
	rows, err := db.Query(query, merchantID, restoranID)
	if err != nil {
		log.Println(err)
		return
	}

	var makanan model.DetailMakanan
	var makanans []model.DetailMakanan
	for rows.Next() {
		if err := rows.Scan(&makanan.MakananID, &makanan.MakananNama, &makanan.MakananHarga); err != nil {
			log.Println(err)
			return
		} else {
			makanans = append(makanans, makanan)
		}
	}

	if len(makanans) < 1 {
		sendResponse(w, 400, "Ayo Buat Makanan")
	} else {
		lihatMakananResponse(w, makanans)
	}
}

package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Tubes/model"
)

func LihatOrder(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	_, userID, _, _ := validateTokenFromCookies(r)

	var userType int
	row := db.QueryRow("SELECT user_type FROM users WHERE user_id=?", userID)
	if err := row.Scan(&userType); err != nil {
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Login as Courrier First")
		return
	}

	if userID <= 0 || userType != 3 {
		sendResponse(w, 400, "Access Denied || Login as Courrier First")
		return
	}

	// Select transaction IDs where user ID matches and status is still "Pending"
	rows, err := db.Query("SELECT transaksi.transaksi_id, transaksi.harga_total, dt.makanan_id, dt.kuantitas FROM transaksi JOIN detail_transaksi dt ON transaksi.transaksi_id = dt.transaksi_id where transaksi.status = ?", 2)
	if err != nil {
		log.Println("Error selecting user orders:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Loop through results and append transaction IDs to a slice
	var transaksi model.DetailTransaksiOrder
	var transaksis []model.DetailTransaksiOrder
	for rows.Next() {
		if err := rows.Scan(&transaksi.TransakasiID, &transaksi.TransakasiHargaTotal, &transaksi.TransaksiDetailOrder.DetailTransaksiMakananID, &transaksi.TransaksiDetailOrder.DetailTransaksiKuantitas); err != nil {
			log.Println(err)
			return
		} else {
			transaksis = append(transaksis, transaksi)
		}
	}

	if len(transaksis) < 1 {
		sendResponse(w, 400, "Belum ada Transaksi")
	} else {
		LihatDetailOrderResponse(w, transaksis)
	}
}

func AmbilOrderan(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	transaksiIDStr := r.Form.Get("transaksi_id")
	transaksiID, err := strconv.Atoi(transaksiIDStr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, userID, _, _ := validateTokenFromCookies(r)

	var userType int
	row := db.QueryRow("SELECT user_type FROM users WHERE user_id=?", userID)
	if err := row.Scan(&userType); err != nil {
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Login as Courrier First")
		return
	}

	if userID <= 0 || userType != 3 {
		sendResponse(w, 400, "Access Denied || Login as Courrier First")
		return
	}

	var memberID int
	row = db.QueryRow("SELECT user_id FROM transaksi WHERE transaksi_id=?", transaksiID)
	if err := row.Scan(&memberID); err != nil {
		log.Println(userID)
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Login as Courrier First")
		return
	}

	var status int
	row = db.QueryRow("SELECT status FROM transaksi WHERE transaksi_id=?", transaksiID)
	if err := row.Scan(&status); err != nil {
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Login as Courrier First")
		return
	}

	var courrierID int
	row = db.QueryRow("SELECT courrier_id FROM courrier WHERE user_id=?", userID)
	if err := row.Scan(&courrierID); err != nil {
		log.Println("Error scanning user count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var count int
	row = db.QueryRow("SELECT COUNT(*) FROM transaksi WHERE transaksi_id=?", transaksiID)
	if err := row.Scan(&count); err != nil {
		log.Println("Error scanning transaksi count:", err)
		sendResponse(w, 500, "Internal Server Error")
		return
	}
	if count == 0 {
		log.Println("Transaksi not found")
		sendResponse(w, 404, "Transaksi not found")
		return
	}
	if status == 2 {
		// Update courrier status
		_, err = db.Exec("UPDATE courrier SET status_pengantaran = '1' WHERE courrier_id = ?", courrierID)
		if err != nil {
			log.Println("Error updating transaksi status:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("UPDATE transaksi SET status = '3' WHERE user_id = ? AND transaksi_id = ?", memberID, transaksiID)
		if err != nil {
			log.Println("Error updating transaksi status:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		sendResponse(w, 200, "Segera Delivery")
	} else {
		sendResponse(w, 400, "Delivery Not Available")
	}
}

func FinishOrderan(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	transaksiIDStr := r.Form.Get("transaksi_id")
	transaksiID, err := strconv.Atoi(transaksiIDStr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, userID, _, _ := validateTokenFromCookies(r)

	var userType int
	row := db.QueryRow("SELECT user_type FROM users WHERE user_id=?", userID)
	if err := row.Scan(&userType); err != nil {
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Login as Courrier First")
		return
	}

	if userID <= 0 || userType != 3 {
		sendResponse(w, 400, "Access Denied || Login as Courrier First")
		return
	}

	var memberID int
	row = db.QueryRow("SELECT user_id FROM transaksi WHERE transaksi_id=?", transaksiID)
	if err := row.Scan(&memberID); err != nil {
		log.Println(userID)
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Login as Courrier First")
		return
	}

	var status int
	row = db.QueryRow("SELECT status FROM transaksi WHERE transaksi_id=?", transaksiID)
	if err := row.Scan(&status); err != nil {
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Login as Courrier First")
		return
	}

	var harga_total float64
	row = db.QueryRow("SELECT harga_total FROM transaksi WHERE transaksi_id=?", transaksiID)
	if err := row.Scan(&harga_total); err != nil {
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Login as Courrier First")
		return
	}

	var untuk_courrier = harga_total * 0.1

	harga_total = harga_total - untuk_courrier

	var restoranID int
	row = db.QueryRow("SELECT restoran_id FROM transaksi WHERE transaksi_id=?", transaksiID)
	if err := row.Scan(&restoranID); err != nil {
		log.Println("Error scanning user count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var courrierID int
	row = db.QueryRow("SELECT courrier_id FROM courrier WHERE user_id=?", userID)
	if err := row.Scan(&courrierID); err != nil {
		log.Println("Error scanning user count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var count int
	row = db.QueryRow("SELECT COUNT(*) FROM transaksi WHERE transaksi_id=?", transaksiID)
	if err := row.Scan(&count); err != nil {
		log.Println("Error scanning transaksi count:", err)
		sendResponse(w, 500, "Internal Server Error")
		return
	}
	if count == 0 {
		log.Println("Transaksi not found")
		sendResponse(w, 404, "Transaksi not found")
		return
	}

	if status == 3 {
		// Update courrier status
		_, err = db.Exec("UPDATE courrier SET status_pengantaran = '0' WHERE courrier_id = ?", courrierID)
		if err != nil {
			log.Println("Error updating transaksi status:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("UPDATE transaksi SET status = '4' WHERE user_id = ? AND transaksi_id = ?", memberID, transaksiID)
		if err != nil {
			log.Println("Error updating transaksi status:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Insert riwayat into database
		_, err = db.Exec("INSERT INTO riwayat_pengantaran (courrier_id, transaksi_id) VALUES (?, ?)", courrierID, transaksiID)
		if err != nil {
			log.Println("Error inserting riwayat:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Update courrier balance
		_, err = db.Exec("UPDATE courrier SET balance = balance + ? WHERE courrier_id = ?", untuk_courrier, courrierID)
		if err != nil {
			log.Println("Error updating transaksi status:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Update restoran balance
		_, err = db.Exec("UPDATE restoran SET balance = balance + ? WHERE restoran_id = ?", harga_total, restoranID)
		if err != nil {
			log.Println("Error updating transaksi status:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		sendResponse(w, 200, "Delivery Finished")
	} else {
		sendResponse(w, 400, "Delivery Canceled")
	}
}

func LihatRiwayatOrderan(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	_, userID, _, _ := validateTokenFromCookies(r)

	var userType int
	row := db.QueryRow("SELECT user_type FROM users WHERE user_id=?", userID)
	if err := row.Scan(&userType); err != nil {
		log.Println("Error scanning user count:", err)
		sendResponse(w, 400, "Login as Courrier First")
		return
	}

	if userID <= 0 || userType != 3 {
		sendResponse(w, 400, "Access Denied || Login as Courrier First")
		return
	}

	query := "select * from riwayat_pengantaran"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var pengantaran model.RiwayatPengantaran
	var pengantarans []model.RiwayatPengantaran
	for rows.Next() {
		if err := rows.Scan(&pengantaran.RiwayatPengantaranID, &pengantaran.RiwayatPengantaranCourrierID, &pengantaran.RiwayatPengantaranTransakasiID); err != nil {
			log.Println(err)
			return
		} else {
			pengantarans = append(pengantarans, pengantaran)
		}
	}

	if len(pengantarans) < 1 {
		sendResponse(w, 400, "Belum ada Riwayat Orderan")
	} else {
		LihatRiwayatPengantaranResponse(w, pengantarans)
	}
}

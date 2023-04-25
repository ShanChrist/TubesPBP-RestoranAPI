package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Tubes/model"
	"github.com/gorilla/mux"
)

func LihatRestoran(w http.ResponseWriter, r *http.Request) {
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
		sendResponse(w, 400, "Blum ada restoran")
	} else {
		lihatRestoranResponse(w, restorans)
	}
}

func LihatSpecificRestoran(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	restoran_id := params["restoran_id"]

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

	query := `SELECT r.restoran_id, r.nama_restoran, r.jam_buka, r.jam_tutup, r.alamat_restoran, r.hari_buka, m.makanan_id, m.nama_makanan, m.harga_makanan
	FROM restoran r
	JOIN makanan m ON r.restoran_id = m.restoran_id
	WHERE r.restoran_id = ?`

	rows, err := db.Query(query, restoran_id)
	if err != nil {
		return
	}
	defer rows.Close()

	restoranMap := make(map[int]model.DetailRestoran)
	for rows.Next() {
		var restoran model.DetailRestoran
		var makanan model.Makanan
		if err := rows.Scan(&restoran.RestoranID, &restoran.RestoranNama, &restoran.RestoranJamBuka, &restoran.RestoranJamTutup, &restoran.RestoranAlamat, &restoran.RestoranHariBuka, &makanan.MakananID, &makanan.MakananNama, &makanan.MakananHarga); err != nil {
			log.Println(err)
			return
		}
		if existingRestoran, ok := restoranMap[restoran.RestoranID]; ok {
			existingRestoran.MakananData = append(existingRestoran.MakananData, makanan)
			restoranMap[restoran.RestoranID] = existingRestoran
		} else {
			restoran.MakananData = append(restoran.MakananData, makanan)
			restoranMap[restoran.RestoranID] = restoran
		}
	}

	var restorans []model.DetailRestoran
	for _, restoran := range restoranMap {
		restorans = append(restorans, restoran)
	}

	if len(restorans) < 1 {
		sendResponse(w, 400, "Tidak ada Menu di restoran ini")
	} else {
		lihatSpecificRestoranResponse(w, restorans)
	}
}

func UserOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	restoran_id := params["restoran_id"]

	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	jumlah_makananStr := r.Form.Get("jumlah_makanan")
	jumlah_makanan, err := strconv.Atoi(jumlah_makananStr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	makananIDStr := r.Form.Get("makanan_id")
	makananID, err := strconv.Atoi(makananIDStr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var makanan model.Makanan
	makanan.MakananID = makananID

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

	// Check if makanan exists for the given restoran_id
	var harga_makanan int
	err = db.QueryRow("SELECT harga_makanan FROM makanan WHERE restoran_id = ? AND makanan_id = ?", restoran_id, makananID).Scan(&harga_makanan)
	if err != nil {
		log.Println("Error checking makanan:", err)
		sendResponse(w, 400, "Tidak ada menu tersebut di restoran kami")
		return
	}
	if harga_makanan == 0 {
		http.Error(w, "Makanan not found", http.StatusNotFound)
		return
	}

	total_harga := jumlah_makanan * harga_makanan

	// Insert data to transaksi table
	_, err = db.Exec("INSERT INTO transaksi (restoran_id, user_id, harga_total, status) VALUES (?, ?, ?, ?)", restoran_id, userID, total_harga, 1)
	if err != nil {
		log.Println("Error inserting transaksi data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// transaksi_id	makanan_id	kuantitas
	var transaksiID int
	row = db.QueryRow("SELECT MAX(transaksi_id) FROM transaksi")

	if err := row.Scan(&transaksiID); err != nil {
		log.Println("Error scanning transaksi count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Insert data to detail_transaksi table
	_, err = db.Exec("INSERT INTO detail_transaksi (transaksi_id, makanan_id, kuantitas) VALUES (?, ?, ?)", transaksiID, makananID, jumlah_makanan)
	if err != nil {
		log.Println("Error inserting transaksi data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sendResponse(w, 200, "Success Ordering, Pay up")
}

func LihatStatusOrder(w http.ResponseWriter, r *http.Request) {
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

	// Select transaction IDs where user ID matches and status is still "Pending"
	rows, err := db.Query("SELECT transaksi.transaksi_id, transaksi.harga_total, transaksi_status.name FROM transaksi JOIN transaksi_status ON transaksi.status = transaksi_status.id WHERE transaksi.user_id=?", userID)
	if err != nil {
		log.Println("Error selecting user orders:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Loop through results and append transaction IDs to a slice
	var transaksi model.DetailMemberTransaksi
	var transaksis []model.DetailMemberTransaksi
	for rows.Next() {
		if err := rows.Scan(&transaksi.TransakasiID, &transaksi.TransakasiHargaTotal, &transaksi.TransakasiStatus.TransaksiStatusName); err != nil {
			log.Println(err)
			return
		} else {
			transaksis = append(transaksis, transaksi)
		}
	}
	if len(transaksis) < 1 {
		sendResponse(w, 400, "Belum ada Transaksi")
	} else {
		LihatStatusOrderResponse(w, transaksis)
	}
}

func Pay_Order(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	transaksiIDStr := r.Form.Get("transaksi_id")
	transaksi_id, err := strconv.Atoi(transaksiIDStr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	total_bayarStr := r.Form.Get("total_bayar")
	totalbayar, err := strconv.Atoi(total_bayarStr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	promo := r.Form.Get("promo")
	point := r.Form.Get("point")

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

	var memberID int
	row = db.QueryRow("SELECT member_id FROM member WHERE user_id=?", userID)

	if err := row.Scan(&memberID); err != nil {
		log.Println("Error scanning user count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var balance int
	err = db.QueryRow("SELECT balance FROM member WHERE member_id = ?", memberID).Scan(&balance)
	if err != nil {
		log.Println("Error checking member:", err)
		return
	}

	var total_harga int
	err = db.QueryRow("SELECT harga_total FROM transaksi WHERE user_id = ? and transaksi_id = ?", userID, transaksi_id).Scan(&total_harga)
	if err != nil {
		log.Println("Error checking transaksi:", err)
		return
	}

	var status int
	err = db.QueryRow("SELECT status FROM transaksi WHERE user_id = ? and transaksi_id = ?", userID, transaksi_id).Scan(&status)
	if err != nil {
		log.Println("Error checking transaksi:", err)
		return
	}

	var point_member int
	err = db.QueryRow("SELECT point FROM member WHERE member_id = ?", memberID).Scan(&point_member)
	if err != nil {
		log.Println("Error checking member:", err)
		return
	}

	var dapat_point = float64(total_harga) * 0.1
	var parsing_point = dapat_point

	if status == 1 {
		var promoValue float32 = 0

		if promo != "" {
			err = db.QueryRow("SELECT persentase FROM promo where kode_promo = ? ", promo).Scan(&promoValue)
			if err != nil {
				sendResponse(w, 400, "Promo not Available")
				return
			}
			if promoValue == 0 {
				sendResponse(w, 400, "Promo not Available")
				return
			}
		}
		total_balance := balance + point_member

		if totalbayar == total_harga {

			if total_balance >= total_harga {
				// Update transaksi status to "Paid"
				_, err = db.Exec("UPDATE transaksi SET status = '2' WHERE user_id = ? AND transaksi_id = ?", userID, transaksi_id)
				if err != nil {
					log.Println("Error updating transaksi status:", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				if point == "use" {
					if total_harga > point_member {
						total_harga = total_harga - point_member
						point_member = 0
					} else {
						point_member = point_member - total_harga
						total_harga = 0
					}
				}

				jumlah_saldo := balance - total_harga
				dapat_point = dapat_point + float64(point_member)
				discount := float32(totalbayar) * promoValue
				jumlah_saldo = jumlah_saldo + int(discount)

				_, err = db.Exec("UPDATE member SET balance = ?, point = ? WHERE member_id = ?", jumlah_saldo, dapat_point, memberID)
				if err != nil {
					log.Println("Error updating balance and point:", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return

				}
				SendEmail(w, r, userID, int(parsing_point))
				sendResponse(w, 200, "Payment Success, Please wait the Driver")
			} else {
				sendResponse(w, 400, "Payment Failed. Balance not enough")
			}
		} else {
			sendResponse(w, 400, "Payment Failed. Total Bayar does not match Total Harga")
			return
		}
	} else {
		sendResponse(w, 400, "You already paid this transaction")
	}
}

package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Tubes/model"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if email == "" || password == "" {
		http.Error(w, "False || Invalid Email atau Password", http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM users WHERE email=? AND password=?", email, password)

	var user model.User
	if err := row.Scan(&user.UserID, &user.UserName, &user.UserFirstName, &user.UserLastName, &user.UserPhoneNumber, &user.UserEmail, &user.UserPassword, &user.UserAddress, &user.UserType); err != nil {
		sendResponse(w, 400, "False || Invalid Email/Password")
	} else {
		generateToken(w, user.UserID, user.UserName, user.UserType)
		sendResponse(w, 200, "True || Login Success")
	}
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)

	sendResponse(w, 200, "Logout")
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user := model.User{
		UserName:        r.Form.Get("username"),
		UserFirstName:   r.Form.Get("first_name"),
		UserLastName:    r.Form.Get("last_name"),
		UserPhoneNumber: r.Form.Get("phone_number"),
		UserEmail:       r.Form.Get("email"),
		UserPassword:    r.Form.Get("password"),
		UserAddress:     r.Form.Get("address"),
	}

	userTypeStr := r.Form.Get("user_type")
	userType, err := strconv.Atoi(userTypeStr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	user.UserType = userType

	if user.UserName == "" || user.UserFirstName == "" || user.UserLastName == "" || user.UserPhoneNumber == "" || user.UserEmail == "" || user.UserPassword == "" || user.UserAddress == "" {
		http.Error(w, "False || Invalid input", http.StatusBadRequest)
		return
	}
	if user.UserType >= 4 {
		sendResponse(w, 400, "1 = Member | 2 = Merchant | 3 = Courrier")
		return
	}

	var count int
	// Check if username already exist
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE username=?", user.UserName)

	if err := row.Scan(&count); err != nil {
		log.Println("Error scanning user count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "False || Username already exists", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	row = db.QueryRow("SELECT COUNT(*) FROM users WHERE email=?", user.UserEmail)
	if err := row.Scan(&count); err != nil {
		log.Println("Error scanning user count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "False || Email already exists", http.StatusBadRequest)
		return
	}
	// insert user into database
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	result, err := tx.Exec("INSERT INTO users (username, first_name, last_name, phone_number, email, password, address, user_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", user.UserName, user.UserFirstName, user.UserLastName, user.UserPhoneNumber, user.UserEmail, user.UserPassword, user.UserAddress, user.UserType)
	if err != nil {
		log.Println("Error inserting user:", err)
		tx.Rollback()
		// Delete the user from the database
		_, err = tx.Exec("DELETE FROM users WHERE username = ?", user.UserName)
		if err != nil {
			log.Println("Error deleting user:", err)
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// check if insert was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 1 {
		// Get the user ID of the newly created user
		row = db.QueryRow("SELECT user_id FROM users WHERE email = ?", user.UserEmail)
		var userID int
		if err := row.Scan(&userID); err != nil {
			log.Println("Error getting user ID:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		switch userType {
		case 1: // Member

			// Insert data to member table
			_, err := db.Exec("INSERT INTO member (user_id, balance, point) VALUES (?, ?, ?)", userID, 0, 0)
			if err != nil {
				log.Println("Error inserting member data:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		case 2: // Merchant
			// Insert data to merchant table
			_, err := db.Exec("INSERT INTO merchant (user_id, balance) VALUES (?, ?)", userID, 0)
			if err != nil {
				log.Println("Error inserting member data:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

		case 3: // Courrier
			// Insert data to courrier table
			_, err := db.Exec("INSERT INTO courrier (user_id, balance, status_pengantaran) VALUES (?, ?, ?)", userID, 0, 0)
			if err != nil {
				log.Println("Error inserting member data:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		default:
			sendResponse(w, 400, "False || Error creating user")
		}
	}
	sendResponse(w, 200, "True || User created successfully")
}

func Topup(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	topUp := r.Form.Get("topup")
	total_topup, err := strconv.Atoi(topUp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, userID, _, _ := validateTokenFromCookies(r)

	var userType int
	row := db.QueryRow("SELECT user_type FROM users WHERE user_id=?", userID)
	if err := row.Scan(&userType); err != nil {
		sendResponse(w, 400, "Login First")
		return
	}

	switch userType {
	case 1:
		_, err = db.Exec("UPDATE member SET balance = balance + ? WHERE user_id = ?", total_topup, userID)
		if err != nil {
			log.Println("Error updating transaksi status:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case 2:
		_, err = db.Exec("UPDATE merchant SET balance = balance + ? WHERE user_id = ?", total_topup, userID)
		if err != nil {
			log.Println("Error updating transaksi status:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case 3:
		_, err = db.Exec("UPDATE courrier SET balance = balance + ? WHERE user_id = ?", total_topup, userID)
		if err != nil {
			log.Println("Error updating transaksi status:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	sendResponse(w, 200, "Topup Success")
}

func CheckBalance(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	_, userID, _, _ := validateTokenFromCookies(r)

	var userType int
	row := db.QueryRow("SELECT user_type FROM users WHERE user_id=?", userID)
	if err := row.Scan(&userType); err != nil {
		sendResponse(w, 400, "Login First")
		return
	}

	if userID <= 0 {
		sendResponse(w, 400, "Login First")
		return
	}

	var balance int
	switch userType {
	case 1:
		err = db.QueryRow("SELECT balance FROM member WHERE user_id = ?", userID).Scan(&balance)
		if err != nil {
			log.Println("Error checking member:", err)
			return
		}
	case 2:
		err = db.QueryRow("SELECT balance FROM merchant WHERE user_id = ?", userID).Scan(&balance)
		if err != nil {
			log.Println("Error checking member:", err)
			return
		}
	case 3:
		err = db.QueryRow("SELECT balance FROM courrier WHERE user_id = ?", userID).Scan(&balance)
		if err != nil {
			log.Println("Error checking member:", err)
			return
		}
	}
	sendResponse(w, 200, "Balance : "+strconv.Itoa(balance))
}

func LihatPromo(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	_, userID, _, _ := validateTokenFromCookies(r)

	if userID < 0 {
		sendResponse(w, 400, "Access Denied || Login First")
		return
	}

	query := "select * from promo"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var promo model.Promo
	var promos []model.Promo
	for rows.Next() {
		if err := rows.Scan(&promo.PromoKode, &promo.PromoDeskripsi, &promo.PromoPersentase); err != nil {
			log.Println(err)
			return
		} else {
			promos = append(promos, promo)
		}
	}

	if len(promos) < 1 {
		sendResponse(w, 400, "Error Array Size Not Correct")
	} else {
		lihatPromoResponse(w, promos)
	}
}

func EditProfile(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	_, userID, _, _ := validateTokenFromCookies(r)

	if userID < 0 {
		sendResponse(w, 400, "Access Denied || Login First")
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	user := model.User{
		UserName:        r.Form.Get("username"),
		UserFirstName:   r.Form.Get("first_name"),
		UserLastName:    r.Form.Get("last_name"),
		UserPhoneNumber: r.Form.Get("phone_number"),
		UserAddress:     r.Form.Get("address"),
	}

	query := "UPDATE users SET username = ?, first_name = ?, last_name = ?, phone_number = ? , address = ? WHERE user_id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println(err)
		return
	}

	result, err := stmt.Exec(user.UserName, user.UserFirstName, user.UserLastName, user.UserPhoneNumber, user.UserAddress, userID)
	if err != nil {
		log.Println(err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return
	}

	if rowsAffected == 0 {
		sendResponse(w, 400, "The id may not exist in the table.")
		return
	}
	sendResponse(w, 200, "Success Updated")
}

func LihatProfile(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	_, userID, _, _ := validateTokenFromCookies(r)

	if userID < 0 {
		sendResponse(w, 400, "Access Denied || Login First")
		return
	}

	query := "select username, first_name, last_name, phone_number, email, address from users where user_id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	var user model.LihatProfile

	rows, err := stmt.Query(userID)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		found = true
		if err := rows.Scan(&user.UserName, &user.UserFirstName, &user.UserLastName, &user.UserPhoneNumber, &user.UserEmail, &user.UserAddress); err != nil {
			log.Println(err)
			return
		}
	}

	if !found {
		sendResponse(w, 400, "Data Not Found")
		return
	}
	lihatProfileResponse(w, user)
}

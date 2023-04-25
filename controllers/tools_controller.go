package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"gopkg.in/gomail.v2"

	"github.com/Tubes/model"
)

func Cron(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	c := gocron.NewScheduler(time.UTC)
	// c.Every(1).Month().Do(SendEmailSubsMonthly)
	c.Every(10).Second().Do(SendEmailMonthly)
	c.StartAsync()
	sendResponse(w, 200, "Success")
}

func SendEmailMonthly() {
	db := connect()
	defer db.Close()

	var kodePromo string
	var persentase float32
	row := db.QueryRow("SELECT kode_promo, persentase FROM promo ORDER BY kode_promo DESC LIMIT 1")
	if err := row.Scan(&kodePromo, &persentase); err != nil {
		log.Println("Error scanning promo data:", err)
		return
	}
	persentaseStr := fmt.Sprintf("%.0f%%", persentase*100)

	query := "SELECT * FROM users"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var user model.User
	var names []string
	var emails []string
	for rows.Next() {
		if err := rows.Scan(&user.UserID, &user.UserName, &user.UserFirstName, &user.UserLastName, &user.UserPhoneNumber, &user.UserEmail, &user.UserPassword, &user.UserAddress, &user.UserType); err != nil {
			log.Println(err)
			return
		} else {
			emails = append(emails, user.UserEmail)
			names = append(names, user.UserFirstName)
		}
	}

	log.Println(emails)
	emailSender := "ganjarhooyah@outlook.com"
	m := gomail.NewDialer("smtp-mail.outlook.com", 587, emailSender, "Hohohooyah")

	var wg sync.WaitGroup
	for i := 0; i < len(emails); i++ {
		wg.Add(1)
		go func(email string, name string) {
			defer wg.Done()

			mail := gomail.NewMessage()
			mail.SetHeader("From", emailSender)
			mail.SetHeader("To", email)
			mail.SetHeader("Subject", "Subscription")
			mail.SetBody("text/html", "Hi "+name+", its your lucky day, you can receive a "+persentaseStr+"% discount on your order when you use the promo code "+kodePromo+" at checkout. <br>Check the link below for more info <br><a href='chess.com'>click link here</a><br> Sincerely, <br> hohohihehooyah")
			if err := m.DialAndSend(mail); err != nil {
				fmt.Println(err)
				return
			}
			log.Println(email)
		}(emails[i], names[i])
	}
	wg.Wait()
}

func SendEmail(w http.ResponseWriter, r *http.Request, userID int, point int) {
	db := connect()
	defer db.Close()

	pointstr := strconv.Itoa(point)

	var email string
	row := db.QueryRow("SELECT email FROM users WHERE user_id=?", userID)
	if err := row.Scan(&email); err != nil {
		log.Println("Error scanning user count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var firstname string
	row = db.QueryRow("SELECT first_name FROM users WHERE user_id=?", userID)
	if err := row.Scan(&firstname); err != nil {
		log.Println("Error scanning user count:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	m := gomail.NewDialer("smtp-mail.outlook.com", 587, "hohohihehooyah@outlook.com", "Hohohooyah")
	mail := gomail.NewMessage()

	mail.SetHeader("From", "hohohihehooyah@outlook.com")
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "Reset Password")
	mail.SetBody("text/html", "Hey, <b>"+firstname+"</b><br> Transaksi berhasil!! Kamu mendapatkan "+pointstr+" point yang dapat kamu pakai dipembelian selanjutnya!!")

	if err := m.DialAndSend(mail); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Tubes/controllers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	//USER
	router.HandleFunc("/login", controllers.UserLogin).Methods("POST")
	router.HandleFunc("/logout", controllers.UserLogout).Methods("POST")
	router.HandleFunc("/register", controllers.UserRegister).Methods("POST")
	router.HandleFunc("/balance", controllers.CheckBalance).Methods("GET")
	router.HandleFunc("/promo", controllers.LihatPromo).Methods("GET")
	router.HandleFunc("/topup", controllers.Topup).Methods("POST")
	router.HandleFunc("/edit_profile", controllers.EditProfile).Methods("PUT")
	router.HandleFunc("/profile", controllers.LihatProfile).Methods("GET")

	//MEMBER
	router.HandleFunc("/restoran", controllers.Authenticate(controllers.LihatRestoran, 1)).Methods("GET")
	router.HandleFunc("/restoran/{restoran_id}", controllers.Authenticate(controllers.LihatSpecificRestoran, 1)).Methods("GET")
	router.HandleFunc("/order", controllers.Authenticate(controllers.LihatStatusOrder, 1)).Methods("GET")
	router.HandleFunc("/order/{restoran_id}", controllers.Authenticate(controllers.UserOrder, 1)).Methods("POST")
	router.HandleFunc("/pay_order", controllers.Authenticate(controllers.Pay_Order, 1)).Methods("POST")

	http.Handle("/", router)
	fmt.Println("Connected to port 7777")
	log.Println("Connected to port 7777")
	log.Fatal(http.ListenAndServe(":7777", router))
}

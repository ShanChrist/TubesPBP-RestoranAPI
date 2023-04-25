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

	http.Handle("/", router)
	fmt.Println("Connected to port 7777")
	log.Println("Connected to port 7777")
	log.Fatal(http.ListenAndServe(":7777", router))
}

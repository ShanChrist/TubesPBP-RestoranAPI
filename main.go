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

	//courrier
	router.HandleFunc("/get_order", controllers.Authenticate(controllers.LihatOrder, 3)).Methods("GET")
	router.HandleFunc("/ambil_orderan", controllers.Authenticate(controllers.AmbilOrderan, 3)).Methods("POST")
	router.HandleFunc("/finish_orderan", controllers.Authenticate(controllers.FinishOrderan, 3)).Methods("POST")
	router.HandleFunc("/riwayat_orderan", controllers.Authenticate(controllers.LihatRiwayatOrderan, 3)).Methods("GET")

	//MERCHANT
	router.HandleFunc("/register_restoran", controllers.Authenticate(controllers.RestoranRegister, 2)).Methods("POST")
	router.HandleFunc("/insert_makanan", controllers.Authenticate(controllers.InsertMakanan, 2)).Methods("POST")
	router.HandleFunc("/makanan", controllers.Authenticate(controllers.LihatMakanan, 2)).Methods("GET")
	router.HandleFunc("/delete_makanan/{makanan_id}", controllers.Authenticate(controllers.DeleteMakanan, 2)).Methods("DELETE")

	//ADMIN
	router.HandleFunc("/insert_promo", controllers.Authenticate(controllers.InsertPromo, 0)).Methods("POST")
	router.HandleFunc("/delete_promo/{kode_promo}", controllers.Authenticate(controllers.DeletePromo, 0)).Methods("DELETE")
	router.HandleFunc("/cron", controllers.Authenticate(controllers.Cron, 0)).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Connected to port 7777")
	log.Println("Connected to port 7777")
	log.Fatal(http.ListenAndServe(":7777", router))
}

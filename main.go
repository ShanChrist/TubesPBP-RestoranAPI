package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	http.Handle("/", router)
	fmt.Println("Connected to port 7777")
	log.Println("Connected to port 7777")
	log.Fatal(http.ListenAndServe(":7777", router))
}

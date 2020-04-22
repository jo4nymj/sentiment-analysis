package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"code.sentiments/services"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/reviews/{id}", services.GetReviews).Methods("GET")

	router.HandleFunc("/products/{name}", services.GetProduct).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

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
	router.HandleFunc("/reviews/analysis:{id}", services.UpdateAnalysisReviews).Methods("POST")

	router.HandleFunc("/products/{name}", services.GetProduct).Methods("GET")
	router.HandleFunc("/products/{name}", services.UpdateProduct).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"code.sentiments/services"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/reviews/{id}", services.ListReviews).Methods("GET")
	router.HandleFunc("/reviews/rating/{id}", services.UpdateRatingReview).Methods("POST")
	router.HandleFunc("/reviews/rating/{id}", services.CreateRatingReview).Methods("PUT")

	router.HandleFunc("/products/{name}", services.GetProduct).Methods("GET")
	router.HandleFunc("/products/{name}", services.UpdateProduct).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

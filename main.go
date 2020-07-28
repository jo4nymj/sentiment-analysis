package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"code.sentiments/services"
)

func main() {

	router := mux.NewRouter()
	router.Use(middleware)

	router.HandleFunc("/reviews/{id}", services.ListReviews).Methods("GET")
	router.HandleFunc("/reviews/rating/{id}", services.UpdateRatingReview).Methods("POST")
	router.HandleFunc("/reviews/rating/{id}", services.CreateRatingReview).Methods("PUT")
	router.HandleFunc("/reviews/rating/", services.CreateRatingReviewCron).Methods("PUT")

	router.HandleFunc("/products/{name}", services.GetProduct).Methods("GET")
	router.HandleFunc("/products/{name}", services.UpdateProduct).Methods("POST")

	router.HandleFunc("/foo", services.GetFoo).Methods("GET")

	log.Fatal(http.ListenAndServe(":5002", router))
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

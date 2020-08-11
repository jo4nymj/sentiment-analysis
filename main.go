package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"code.sentiments/config"
	"code.sentiments/services"
	"code.sentiments/utilities"
)

func main() {

	router := mux.NewRouter()
	router.Use(utilities.Middleware)
	if config.IsProduction {
		router.Use(utilities.AuthMiddleware)
	}

	router.HandleFunc("/reviews/{id}", services.ListReviews).Methods("GET")
	router.HandleFunc("/reviews/rating/{id}:create", services.CreateRatingReview).Methods("POST")
	router.HandleFunc("/reviews/rating/{id}:update", services.UpdateRatingReview).Methods("POST")

	router.HandleFunc("/products/{name}", services.GetProduct).Methods("GET")
	router.HandleFunc("/products/{name}", services.UpdateProduct).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT"},
	})

	log.Fatal(http.ListenAndServe(":5002", c.Handler(router)))
}

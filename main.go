package main

import (
	"log"
	"net/http"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"code.sentiments/config"
	"code.sentiments/services"
)

var key = ""

func main() {

	router := mux.NewRouter()
	router.Use(middleware)
	if config.IsProduction {
		key = config.MustGetenv("KEY")
		router.Use(authMiddleware)
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

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	if len(key) == 0 {
		log.Fatal("HTTP server unable to start, expected an APP_KEY for JWT auth")
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	return jwtMiddleware.Handler(next)
}

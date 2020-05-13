package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"code.sentiments/config"
	"code.sentiments/repository"
)

func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := config.GetMySQLDB()
	if err != nil {
		fmt.Println(err)
	}

	productModel := repository.ProductModel{
		Db: db,
	}

	product, err := productModel.GetProduct(mux.Vars(r)["name"])
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(product)

}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := config.GetMySQLDB()
	if err != nil {
		fmt.Println(err)
	}

	productModel := repository.ProductModel{
		Db: db,
	}
	product, err := productModel.GetProduct(mux.Vars(r)["name"])
	if err != nil {
		fmt.Println(err)
	}

	reviewModel := repository.ReviewModel{
		Db: db,
	}
	reviews, err := reviewModel.ListReviews(product.ID)
	if err != nil {
		fmt.Println(err)
	}

	var total float32
	var n float32
	for _, review := range reviews {
		log.Info("La puntuación de la review es: ", review.Score)
		total = total + review.Score
		n = n + 1
	}

	product.Average_Rating = total / n

	json.NewEncoder(w).Encode(product)
}

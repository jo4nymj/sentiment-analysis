package services

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"code.sentiments/config"
	"code.sentiments/logger"
	"code.sentiments/repository"
)

func GetProduct(w http.ResponseWriter, r *http.Request) {
	productModel := repository.ProductModel{
		Db: config.Instance,
	}

	product, err := productModel.GetProduct(mux.Vars(r)["name"])
	if err != nil {
		logger.Error("Failed to retrieve the product from database, ", err)
		log.Errorf("Failed to retrieve the product from database, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Something bad happened!"))
	}
	json.NewEncoder(w).Encode(product)

}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productModel := repository.ProductModel{
		Db: config.Instance,
	}
	product, err := productModel.GetProduct(mux.Vars(r)["name"])
	if err != nil {
		logger.Error("Failed to retrieve the product from database, ", err)
		log.Errorf("Failed to retrieve the product from database, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Something bad happened!"))
	}

	reviewModel := repository.ReviewModel{
		Db: config.Instance,
	}
	reviews, err := reviewModel.ListReviews(product.ID)
	if err != nil {
		logger.Error("Failed to retrieve the review from database, ", err)
		log.Errorf("Failed to retrieve the review from database, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Something bad happened!"))
	}

	var total float32
	var n float32
	for _, review := range reviews {
		total = total + review.Rating
		n = n + 1
	}
	product.Average_Rating = total / n

	if err := productModel.UpdateProduct(product); err != nil {
		logger.Error("Failed to update the product ", err)
		log.Errorf("Failed to update the product, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Something bad happened!"))
	}

	json.NewEncoder(w).Encode(product)
}

func GetFoo(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Hola Mundo")
}

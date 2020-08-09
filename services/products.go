package services

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"code.sentiments/config"
	"code.sentiments/logger"
	"code.sentiments/repository"
	"code.sentiments/utilities"
)

func GetProduct(w http.ResponseWriter, r *http.Request) {
	productModel := repository.ProductModel{
		Db: config.Instance,
	}

	product, err := productModel.GetProduct(mux.Vars(r)["name"])
	if err != nil {
		logger.Error("Failed to retrieve the product from database, %v", err)
		utilities.StatusBadRequest(w, r, "")
	}
	json.NewEncoder(w).Encode(product)

}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productModel := repository.ProductModel{
		Db: config.Instance,
	}
	product, err := productModel.GetProduct(mux.Vars(r)["name"])
	if err != nil {
		logger.Error("Failed to retrieve the product from database, %v", err)
		utilities.StatusBadRequest(w, r, "")
	}

	reviewModel := repository.ReviewModel{
		Db: config.Instance,
	}
	reviews, err := reviewModel.ListReviews(product.ID)
	if err != nil {
		logger.Error("Failed to retrieve the review from database, %v", err)
		utilities.StatusBadRequest(w, r, "")
	}

	var total float32
	var n float32
	for _, review := range reviews {
		total = total + review.Rating
		n = n + 1
	}
	product.Average_Rating = total / n

	if err := productModel.UpdateProduct(product); err != nil {
		logger.Error("Failed to update the product %v", err)
		utilities.StatusBadRequest(w, r, "")
	}

	json.NewEncoder(w).Encode(product)
}

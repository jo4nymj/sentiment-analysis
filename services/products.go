package services

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"

	"code.sentiments/config"
	"code.sentiments/repository"
)

func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	productModel := repository.ProductModel{
		Db: config.Instance,
	}

	product, err := productModel.GetProduct(mux.Vars(r)["name"])
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(product)

}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	productModel := repository.ProductModel{
		Db: config.Instance,
	}
	product, err := productModel.GetProduct(mux.Vars(r)["name"])
	if err != nil {
		fmt.Println(err)
	}

	reviewModel := repository.ReviewModel{
		Db: config.Instance,
	}
	reviews, err := reviewModel.ListReviews(product.ID)
	if err != nil {
		fmt.Println(err)
	}

	var total float32
	var n float32
	for _, review := range reviews {
		total = total + review.Rating
		n = n + 1
	}

	product.Average_Rating = total / n

	if err := productModel.UpdateProduct(product); err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(product)
}

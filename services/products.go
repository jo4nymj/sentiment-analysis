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

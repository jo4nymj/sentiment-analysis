package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"code.sentiments/config"
	"code.sentiments/repository"
)

func GetReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := config.GetMySQLDB()
	if err != nil {
		fmt.Println(err)
	}

	reviewModel := repository.ReviewModel{
		Db: db,
	}

	reviewID, err := strconv.ParseInt((mux.Vars(r)["id"]), 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	reviews, err := reviewModel.GetReviews(reviewID)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(reviews)

}

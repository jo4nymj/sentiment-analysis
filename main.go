package main

import (
	"fmt"

	"code.sentiments/config"
	"code.sentiments/models"
)

func main() {
	db, err := config.GetMySQLDB()
	if err != nil {
		fmt.Println(err)
	}

	reviewModel := models.ReviewModel{
		Db: db,
	}

	reviews, err := reviewModel.GetReviews(13)
	if err != nil {
		fmt.Println(err)
	}
	for _, review := range reviews {
		fmt.Println("ID: ", review.ID)
		fmt.Println("Author: ", review.Author)
		fmt.Println("Content: ", review.Content)
	}

	productModel := models.ProductModel{
		Db: db,
	}
	product, err := productModel.GetProduct("NVIDIA")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ID: ", product.ID)
	fmt.Println("Name: ", product.Name)
	fmt.Println("Average_Rating: ", product.Average_Rating)
}

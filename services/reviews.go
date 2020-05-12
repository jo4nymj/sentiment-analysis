package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	//"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"

	language "cloud.google.com/go/language/apiv1"
	"code.sentiments/config"
	"code.sentiments/repository"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
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

func UpdateAnalysisReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := config.GetMySQLDB()
	if err != nil {
		log.Fatal(err)
	}

	reviewModel := repository.ReviewModel{
		Db: db,
	}

	reviewID, err := strconv.ParseInt((mux.Vars(r)["id"]), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	review, err := reviewModel.GetReview(reviewID)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//text := "I'm happy"

	v, err := analyzeSentiment(ctx, client, review.Content)
	if err != nil {
		log.Fatal(err)
	}
	sentiment := v.GetDocumentSentiment()
	review.Score = sentiment.GetScore()
	//json.NewEncoder(w).Encode(proto.MarshalText(w, v))
	json.NewEncoder(w).Encode(review)
}

func analyzeSentiment(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeSentimentResponse, error) {
	return client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}

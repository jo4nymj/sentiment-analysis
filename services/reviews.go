package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	//"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"libs.altipla.consulting/errors"

	language "cloud.google.com/go/language/apiv1"
	"code.sentiments/config"
	"code.sentiments/repository"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

func ListReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reviewModel := repository.ReviewModel{
		Db: config.Instance,
	}

	reviewID, err := strconv.ParseInt((mux.Vars(r)["id"]), 10, 64)
	if err != nil {
		errors.Trace(err)
		return
	}
	reviews, err := reviewModel.ListReviews(reviewID)
	if err != nil {
		errors.Trace(err)
		return
	}

	json.NewEncoder(w).Encode(reviews)

}

func UpdateAnalysisReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reviewModel := repository.ReviewModel{
		Db: config.Instance,
	}

	reviewID, err := strconv.ParseInt((mux.Vars(r)["id"]), 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	review, err := reviewModel.GetReview(reviewID)
	if err != nil {
		fmt.Println("Fallo al hacer get de la review: ", err)
	}

	ctx := context.Background()
	client, err := language.NewClient(ctx)
	if err != nil {
		fmt.Println("Fallo al crear el cliente: ", err)
	}

	v, err := analyzeSentiment(ctx, client, review.Content)
	if err != nil {
		fmt.Println("Fallo en el análisis de la review: ", err)
	}
	sentiment := v.GetDocumentSentiment()
	review.Rating = sentiment.GetScore()
	review.Magnitude = sentiment.GetMagnitude()

	//json.NewEncoder(w).Encode(proto.MarshalText(w, v))

	if err := reviewModel.UpdateReview(review); err != nil {
		fmt.Println("Fallo en el update: ", err)
	}
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

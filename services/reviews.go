package services

import (
	"context"
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"net/http"
	"strconv"

	//"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"libs.altipla.consulting/errors"

	language "cloud.google.com/go/language/apiv1"
	"code.sentiments/config"
	"code.sentiments/models"
	"code.sentiments/repository"
	//log "github.com/sirupsen/logrus"
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

func CreateRatingReview(w http.ResponseWriter, r *http.Request) {
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

	if err := analyze(&review); err != nil {
		fmt.Println("Error al llamar a la API de Google")
	}

	if err := reviewModel.CreateReview(review); err != nil {
		fmt.Println("Fallo en el create: ", err)
	}

	json.NewEncoder(w).Encode(review)
}

func UpdateRatingReview(w http.ResponseWriter, r *http.Request) {
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

	if err := analyze(&review); err != nil {
		fmt.Println("Error al llamar a la API de Google")
	}

	if err := reviewModel.UpdateReview(review); err != nil {
		fmt.Println("Fallo en el update: ", err)
	}

	json.NewEncoder(w).Encode(review)
}

func analyze(review *models.Review) error {
	ctx := context.Background()
	client, err := language.NewClient(ctx)
	if err != nil {
		return err
	}
	v, err := analyzeSentiment(ctx, client, review.Content)
	if err != nil {
		return err
	}
	sentiment := v.GetDocumentSentiment()
	review.Rating = assignStars(sentiment.GetScore())
	review.Magnitude = sentiment.GetMagnitude()

	return nil
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

func assignStars(rating float32) float32 {
	switch {
	case rating >= -1.0 && rating < -0.7:
		return 1.0
	case rating >= -0.7 && rating < -0.3:
		return 2.0
	case rating >= -0.3 && rating < 0.3:
		return 3.0
	case rating >= 0.3 && rating < 0.7:
		return 4.0
	case rating >= 0.7 && rating <= 1.0:
		return 5.0
	default:
		return 0.0
	}
}

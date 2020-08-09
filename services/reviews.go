package services

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	language "cloud.google.com/go/language/apiv1"
	"github.com/gorilla/mux"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"

	"code.sentiments/config"
	"code.sentiments/logger"
	"code.sentiments/models"
	"code.sentiments/repository"
	"code.sentiments/utilities"
)

func ListReviews(w http.ResponseWriter, r *http.Request) {
	reviewModel := repository.ReviewModel{
		Db: config.Instance,
	}

	productID, err := strconv.ParseInt((mux.Vars(r)["id"]), 10, 64)
	if err != nil {
		logger.Error("Failed to parse the product ID, %v", err)
		utilities.StatusInternalServerError(w, r)
	}

	reviews, err := reviewModel.ListReviews(productID)
	if err != nil {
		logger.Error("Failed to list the reviews, %v", err)
		utilities.StatusBadRequest(w, r, "")
	}

	json.NewEncoder(w).Encode(reviews)
}

func CreateRatingReview(w http.ResponseWriter, r *http.Request) {
	reviewModel := repository.ReviewModel{
		Db: config.Instance,
	}

	reviewID, err := strconv.ParseInt((mux.Vars(r)["id"]), 10, 64)
	if err != nil {
		logger.Error("Failed to parse the ID of the review, %v", err)
		utilities.StatusInternalServerError(w, r)
	}

	exists, err := reviewModel.Exists(reviewID)
	if err != nil {
		logger.Error("Failed to retrieve the review from the database, %v", err)
		utilities.StatusBadRequest(w, r, "")
	}

	if exists {
		utilities.StatusBadRequest(w, r, "Review already exists")
		return
	}

	review, err := reviewModel.GetReview(reviewID)
	if err != nil {
		logger.Error("Failed to retrieve the review from the database, %v", err)
		utilities.StatusBadRequest(w, r, "")
	}

	if err := analyze(&review); err != nil {
		logger.Error("Failed to analyze the review, %v", err)
		utilities.StatusInternalServerError(w, r)
	}

	if err := reviewModel.CreateReview(review); err != nil {
		logger.Error("Failed to create review: %v", err)
		utilities.StatusBadRequest(w, r, "")
	}

	json.NewEncoder(w).Encode(review)
}

func UpdateRatingReview(w http.ResponseWriter, r *http.Request) {
	reviewModel := repository.ReviewModel{
		Db: config.Instance,
	}

	reviewID, err := strconv.ParseInt((mux.Vars(r)["id"]), 10, 64)
	if err != nil {
		logger.Error("Failed to parse the ID of the review, %v", err)
		utilities.StatusInternalServerError(w, r)
	}

	review, err := reviewModel.GetReview(reviewID)
	if err != nil {
		logger.Error("Failed to retrieve the review from the database, %v", err)
		utilities.StatusBadRequest(w, r, "")
	}

	if err := analyze(&review); err != nil {
		logger.Error("Failed to analyze the review, %v", err)
		utilities.StatusInternalServerError(w, r)
	}

	if err := reviewModel.UpdateReview(review); err != nil {
		logger.Error("Failed to update review: %v", err)
		utilities.StatusBadRequest(w, r, "")
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

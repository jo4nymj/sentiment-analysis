package services

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	//"github.com/golang/protobuf/proto"
	language "cloud.google.com/go/language/apiv1"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"

	"code.sentiments/config"
	"code.sentiments/logger"
	"code.sentiments/models"
	"code.sentiments/repository"
)

func ListReviews(w http.ResponseWriter, r *http.Request) {
	reviewModel := repository.ReviewModel{
		Db: config.Instance,
	}

	productID, err := strconv.ParseInt((mux.Vars(r)["id"]), 10, 64)
	if err != nil {
		logger.Error("Failed to parse the product ID, ", err)
		log.Errorf("Failed to parse the product ID, %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Something bad happened!"))
	}
	reviews, err := reviewModel.ListReviews(productID)
	if err != nil {
		logger.Error("Failed to list the reviews, ", err)
		log.Errorf("Failed to list the reviews, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Something bad happened!"))
	}

	json.NewEncoder(w).Encode(reviews)
}

func CreateRatingReview(w http.ResponseWriter, r *http.Request) {
	reviewModel := repository.ReviewModel{
		Db: config.Instance,
	}

	reviewID, err := strconv.ParseInt((mux.Vars(r)["id"]), 10, 64)
	if err != nil {
		logger.Error("Failed to parse the ID of the review, ", err)
		log.Errorf("Failed to parse the ID of the review, %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Something bad happened!"))
	}

	review, err := reviewModel.GetReview(reviewID)
	if err != nil {
		logger.Error("Failed to retrieve the review from the database, ", err)
		log.Errorf("Failed to retrieve the review from the database: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Something bad happened!"))
	}

	if err := analyze(&review); err != nil {
		logger.Error("Failed to analyze the review, ", err)
		log.Errorf("Failed to analyze the review: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Something bad happened!"))
	}

	if err := reviewModel.CreateReview(review); err != nil {
		logger.Error("Failed to create review: %v", err)
		log.Errorf("Failed to create review: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Something bad happened!"))
	}

	json.NewEncoder(w).Encode(review)
}

func UpdateRatingReview(w http.ResponseWriter, r *http.Request) {
	reviewModel := repository.ReviewModel{
		Db: config.Instance,
	}

	reviewID, err := strconv.ParseInt((mux.Vars(r)["id"]), 10, 64)
	if err != nil {
		logger.Error("Failed to parse the ID of the review, ", err)
		log.Errorf("Failed to parse the ID of the review, %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Something bad happened!"))
	}

	review, err := reviewModel.GetReview(reviewID)
	if err != nil {
		logger.Error("Failed to retrieve the review from the database, ", err)
		log.Errorf("Failed to retrieve the review from the database: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Something bad happened!"))
	}

	if err := analyze(&review); err != nil {
		logger.Error("Failed to analyze the review, ", err)
		log.Errorf("Failed to analyze the review: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Something bad happened!"))
	}

	if err := reviewModel.UpdateReview(review); err != nil {
		logger.Error("Failed to update review: %v", err)
		log.Errorf("Failed to update review: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Something bad happened!"))
	}

	json.NewEncoder(w).Encode(review)
}

func CreateRatingReviewCron(w http.ResponseWriter, r *http.Request) {
	productModel := repository.ProductModel{
		Db: config.Instance,
	}

	products, err := productModel.ListProducts()
	if err != nil {
		logger.Error("Failed to list products, ", err)
		log.Errorf("Failed to list products, %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Something bad happened!"))
	}

	reviewModel := repository.ReviewModel{
		Db: config.Instance,
	}

	for _, product := range products {
		reviews, err := reviewModel.ListReviews(product.ID)
		if err != nil {
			logger.Error("Failed to list reviews, ", err)
			log.Errorf("Failed to list reviews, %v", err)

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Something bad happened!"))
		}

		analyzedReviews, err := reviewModel.ListAnalyzedReviews(product.ID)
		if err != nil {
			logger.Error("Failed to list analyzed reviews, ", err)
			log.Errorf("Failed to list analyzed reviews, %v", err)

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Something bad happened!"))
		}

		for _, review := range reviews {
			broke := false
			for _, analyzedReview := range analyzedReviews {
				if review.ID == analyzedReview.ID {
					broke = true
					break
				}
			}
			if broke {
				continue
			}

			if err := analyze(&review); err != nil {
				logger.Error("Failed to analyze the review, ", err)
				log.Errorf("Failed to analyze the review: %v", err)

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 Something bad happened!"))
			}

			if err := reviewModel.CreateReview(review); err != nil {
				logger.Error("Failed to create review: %v", err)
				log.Errorf("Failed to create review: %v", err)

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 Something bad happened!"))
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 All reviews have been evaluated succesfully"))
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

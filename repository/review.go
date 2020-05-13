package repository

import (
	//"database/sql"

	"code.sentiments/config"
	"code.sentiments/models"
)

type ReviewModel struct {
	//Db *sql.DB
	Db *config.Connection
}

func (r ReviewModel) ListReviews(ID int64) ([]models.Review, error) {
	rows, err := r.Db.Conn.Query(`SELECT comments.comment_ID, comment_author, comment_content, meta_value
	 FROM wp_comments comments INNER JOIN wp_commentmeta meta ON comments.comment_ID = meta.comment_ID
	 WHERE comment_type = 'review' AND meta.meta_key = 'rating' AND comment_post_ID = ?`, +ID)
	if err != nil {
		return nil, err
	}

	reviews := []models.Review{}
	review := models.Review{}
	for rows.Next() {
		if err := rows.Scan(&review.ID, &review.Author, &review.Content, &review.Score); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r ReviewModel) GetReview(ID int64) (models.Review, error) {
	review := models.Review{}
	rows, err := r.Db.Conn.Query("SELECT meta_value FROM wp_commentmeta WHERE meta_key = 'rating' AND comment_ID = ?", +ID)
	if err != nil {
		return review, err
	}
	for rows.Next() {
		if err := rows.Scan(&review.ID, &review.Author, &review.Content); err != nil {
			return review, err
		}
	}
	return review, nil
}

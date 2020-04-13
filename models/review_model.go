package models

import (
	"database/sql"

	"code.sentiments/entities"
)

type ReviewModel struct {
	Db *sql.DB
}

func (r ReviewModel) GetReviews(ID int64) ([]entities.Review, error) {
	rows, err := r.Db.Query("SELECT comment_ID, comment_author, comment_content FROM wp_comments WHERE comment_type = 'review' AND comment_post_ID = ?", +ID)
	if err != nil {
		return nil, err
	}

	reviews := []entities.Review{}
	review := entities.Review{}
	for rows.Next() {
		if err := rows.Scan(&review.ID, &review.Author, &review.Content); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

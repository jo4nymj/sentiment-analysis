package repository

import (
	"code.sentiments/config"
	"code.sentiments/models"
)

type ReviewModel struct {
	Db *config.Connection
}

func (r ReviewModel) ListReviews(ID int64) ([]models.Review, error) {
	rows, err := r.Db.Conn.Query(`SELECT comments.comment_ID, comment_author, comment_content, meta_value
	 FROM wp_comments comments INNER JOIN wp_commentmeta meta ON comments.comment_ID = meta.comment_ID
	 WHERE comment_type = 'review' AND meta.meta_key = 'verified' AND comment_post_ID = ?`, +ID)
	if err != nil {
		return nil, err
	}

	reviews := []models.Review{}
	review := models.Review{}
	for rows.Next() {
		if err := rows.Scan(&review.ID, &review.Author, &review.Content, &review.Rating); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r ReviewModel) ListAnalyzedReviews(ID int64) ([]models.Review, error) {
	rows, err := r.Db.Conn.Query(`SELECT comments.comment_ID
	 FROM wp_comments comments INNER JOIN wp_commentmeta meta ON comments.comment_ID = meta.comment_ID
	 WHERE comment_type = 'review' AND meta.meta_key = 'rating' AND comment_post_ID = ?`, +ID)
	if err != nil {
		return nil, err
	}

	reviews := []models.Review{}
	review := models.Review{}
	for rows.Next() {
		if err := rows.Scan(&review.ID); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r ReviewModel) GetReview(ID int64) (models.Review, error) {
	review := models.Review{}
	rows, err := r.Db.Conn.Query("SELECT comment_ID, comment_author, comment_content FROM wp_comments WHERE comment_type = 'review' AND comment_ID = ?", +ID)
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

func (r ReviewModel) GetRating(ID int64) (models.Review, error) {
	review := models.Review{}
	rows, err := r.Db.Conn.Query(`SELECT comments.comment_ID, meta_value, meta_magnitude
	FROM wp_comments comments INNER JOIN wp_commentmeta meta ON comments.comment_ID = meta.comment_ID
	WHERE comment_type = 'review' AND meta.meta_key = 'rating' AND comments.comment_ID = ?`, +ID)
	if err != nil {
		return review, err
	}
	for rows.Next() {
		if err := rows.Scan(&review.ID, &review.Rating, &review.Magnitude); err != nil {
			return review, err
		}
	}

	return review, nil
}

func (r ReviewModel) UpdateReview(review models.Review) error {
	stmt, err := r.Db.Conn.Prepare(`UPDATE  wp_commentmeta SET meta_value = ?, meta_magnitude = ? 
		WHERE comment_id = ? AND meta_key = 'rating'`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(review.Rating, review.Magnitude, review.ID); err != nil {
		return err
	}
	return nil
}

func (r ReviewModel) CreateReview(review models.Review) error {
	stmt, err := r.Db.Conn.Prepare(`INSERT INTO wp_commentmeta (comment_id, meta_key, meta_value, meta_magnitude)
		VALUES (?, 'rating', ?, ?)`)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(review.ID, review.Rating, review.Magnitude); err != nil {
		return err
	}

	return nil
}

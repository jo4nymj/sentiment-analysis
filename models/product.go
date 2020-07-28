package models

type Product struct {
	ID             int64   `json:"ID"`
	Name           string  `json:"post_title"`
	Average_Rating float32 `json:"average_rating"`
}

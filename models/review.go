package models

type Review struct {
	ID        int64   `json:"ID"`
	Author    string  `json:"comment_author"`
	Content   string  `json:"comment_content"`
	Rating    float32 `json:"meta_value"`
	Magnitude float32 `json:"meta_magnitude"`
}

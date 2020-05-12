package models

type Review struct {
	ID      int64   `json:"ID"`
	Author  string  `json:"Author"`
	Content string  `json:"Content"`
	Rating  float32 `json:"Rating"`
	Score   float32 `json:"Score"`
}

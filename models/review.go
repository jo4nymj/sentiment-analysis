package models

type Review struct {
	ID      int64   `json:"ID"`
	Author  string  `json:"Author"`
	Content string  `json:"Content"`
	Rating  float64 `json:"Rating"`
	Score   float64 `json:"Score"`
}

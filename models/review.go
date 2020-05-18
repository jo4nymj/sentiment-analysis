package models

type Review struct {
	ID        int64   `json:"ID"`
	Author    string  `json:"Author"`
	Content   string  `json:"Content"`
	Rating    float32 `json:"Rating"`
	Magnitude float32 `json:"Magnitude"`
}

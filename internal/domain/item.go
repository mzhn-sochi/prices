package domain

type Item struct {
	Product     string  `json:"product"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Measure     struct {
		Amount float64 `json:"amount"`
		Unit   string  `json:"unit"`
	} `json:"measure"`
}

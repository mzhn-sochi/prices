package domain

type Item struct {
	Product     string  `json:"product"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Measure     struct {
		Amount uint   `json:"amount"`
		Unit   string `json:"unit"`
	}
}

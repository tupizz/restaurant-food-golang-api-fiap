package dto

type ProductOutput struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Images      []string `json:"images"`
}

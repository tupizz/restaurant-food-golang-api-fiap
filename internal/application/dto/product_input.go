package dto

type ProductInput struct {
	Name        string   `json:"name"`
	Price       string   `json:"price"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Images      []string `json:"images"`
}

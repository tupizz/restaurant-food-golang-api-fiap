package entities

import "time"

type ProductImage struct {
	ID        int
	ImageURL  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductCategory struct {
	ID        int
	Name      string
	Handle    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Category    ProductCategory
	Images      []ProductImage
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

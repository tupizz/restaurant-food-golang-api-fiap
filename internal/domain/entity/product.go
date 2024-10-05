package entity

type ProductImage struct {
	ID       int
	ImageURL string
}

type ProductCategory struct {
	ID   int
	Name string
}

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Category    ProductCategory
	Images      []ProductImage
}

package models

type Product struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Price      int       `json:"price"`
	Stock      int       `json:"stock"`
	CategoryID int       `json:"category_id"`
	Category   *Category `json:"category,omitempty"`
}

type ProductResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type ProductCategoryResponse struct {
	Name string `json:"name"`
}

type ProductDetailResponse struct {
	ID       int                      `json:"id"`
	Name     string                   `json:"name"`
	Price    int                      `json:"price"`
	Stock    int                      `json:"stock"`
	Category *ProductCategoryResponse `json:"category,omitempty"`
}

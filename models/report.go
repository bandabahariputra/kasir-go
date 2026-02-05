package models

type BestSellingProduct struct {
	Name         string `json:"name"`
	QuantitySold int    `json:"quantity_sold"`
}

type TodayReport struct {
	TotalRevenue       int                 `json:"total_revenue"`
	TotalTransaction   int                 `json:"total_transaction"`
	BestSellingProduct *BestSellingProduct `json:"best_selling_product,omitempty"`
}

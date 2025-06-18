package models

type Product struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Price       float64  `json:"price"`
    Description *string  `json:"description,omitempty"`
    Inventory   int      `json:"inventory"`
}

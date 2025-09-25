package models

type Product struct {
    ID          string  `json:"id" gorm:"primaryKey"`
    Name        string  `json:"name"`
    Price       float64 `json:"price"`
    Description string  `json:"description"`
    Inventory   int     `json:"inventory"`
    Available   bool    `json:"available"`
}

func (Product) IsEntity() {}
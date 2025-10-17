package models
// todo change created_at to something simlar to "CreatedAt: s.CreatedAt.Format(time.RFC3339)," see job story story_mapper for example
type Order struct {
	 ID        string  `json:"id" gorm:"primarykey"`  
	UserID     string  `json:"userId"`
	Products []Product `gorm:"many2many:order_products;"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"totalPrice"`
	Status     string  `json:"status"`
	CreatedAt  Time    `json:"createdAt"`
	
	
}


type CreateOrderInput struct {
	UserID     string  `json:"userId"`
	ProductIDs  []string  `json:"productIds"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"totalPrice"`
	Status     string  `json:"status"`
	CreatedAt  Time    `json:"createdAt" gorm:"autoCreateTime"`
}

type UpdateOrderInput struct {
	OrderID     string   `json:"orderId"`
	Quantity    *int     `json:"quantity"`
	TotalPrice  *float64 `json:"totalPrice"`
	Status      *string  `json:"status"`
}

type DeleteOrderInput struct {
	OrderID     string   `json:"orderId"`
	UserID	    string   `json:"userId"`
}

type SetOrderStatusInput struct {
	OrderID     string   `json:"orderId"`
	Status      *string  `json:"status"`
}

type ChangeOrderQuantityInput struct {
	OrderID     string   `json:"orderId"`
	Quantity   int     `json:"quantity"`
}

func (Order) IsEntity() {}






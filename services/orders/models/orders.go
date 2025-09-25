package models

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
	CreatedAt  Time    `json:"createdAt"` 
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






package resolvers

import (
	"github.com/tagaertner/e-commerce-graphql/services/orders/services"
	"gorm.io/gorm"
)

type Resolver struct {
	OrderService *services.OrderService
}

func NewResolver(db *gorm.DB) *Resolver {
	return &Resolver{
		OrderService: services.NewOrderService(db),
	}
}

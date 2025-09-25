package resolvers

import (
	"github.com/tagaertner/e-commerce-graphql/services/orders/generated"
	"github.com/tagaertner/e-commerce-graphql/services/orders/models"
)

func ToGraphQLOrder(o *models.Order) *generated.Order {
    var gqlProducts []*generated.Product
    for _, p := range o.Products {
        gqlProducts = append(gqlProducts, &generated.Product{
            ID:   p.ID,
            Name: p.Name,   
        })
    }
    return &generated.Order{
        ID:         o.ID,
        UserID:     o.UserID,
        Products:   gqlProducts,
        Quantity:   o.Quantity,
        TotalPrice: o.TotalPrice,
        Status:     o.Status,
        CreatedAt:  o.CreatedAt,
    }
}

func ToGraphQLOrders(orders []*models.Order) []*generated.Order {
	var graphQLOrders []*generated.Order
	for _, order := range orders {
		graphQLOrders = append(graphQLOrders, ToGraphQLOrder(order))
	}
	return graphQLOrders
}


func ToGraphQLUser(u *models.User) *generated.User {
    return &generated.User{
        ID: u.ID,
    }
}

func ToGraphQLProduct(p *models.Product) *generated.Product {
	return &generated.Product{
		ID: p.ID,
	}
}

package resolvers

import (
    "github.com/tagaertner/e-commerce-graphql/services/orders/models"
)

func ToGraphQLOrder(o *models.Order) *models.Order {
    return o
}

func ToGraphQLOrders(orders []*models.Order) []*models.Order {
    return orders
}

func ToGraphQLUser(u *models.User) *models.User {
    return u
}

func ToGraphQLProduct(p *models.Product) *models.Product {
    return p
}
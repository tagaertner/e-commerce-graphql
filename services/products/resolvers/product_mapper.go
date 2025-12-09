package resolvers

import (
	// "github.com/tagaertner/e-commerce-graphql/services/products/generated"
	"github.com/tagaertner/e-commerce-graphql/services/products/models"
)

func ToGraphQLProduct(p *models.Product) *models.Product {
	return p
}

func ToGraphQLProductList(products []*models.Product) []*models.Product {
	var gqlProducts []*models.Product
	for _, p := range products {
		gqlProducts = append(gqlProducts, ToGraphQLProduct(p))
	}
	return gqlProducts
}

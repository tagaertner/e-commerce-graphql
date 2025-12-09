package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"github.com/google/uuid"

	// "github.com/tagaertner/e-commerce-graphql/services/products/generated"
	"github.com/tagaertner/e-commerce-graphql/services/products/models"
	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

func (s *ProductService) GetAllProducts() ([]*models.Product, error) {
	var products []*models.Product
	if err := s.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetProductByID(id string) (*models.Product, error) {
	var product models.Product
	if err := s.db.First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *ProductService) CreateProduct(ctx context.Context,  name string, price float64, description string, inventory int ) (*models.Product, error){

	if strings.TrimSpace(name) == ""{
		return nil, fmt.Errorf("invalid product name: missing or invalid field")
	}
	if price <= 0 {
		return nil, fmt.Errorf("price must be greater than zero")
	}
	if inventory < 0 {
		return nil, fmt.Errorf(" inventory cannot be negative")
	}

	product := &models.Product{
		ID:    uuid.NewString(),
		Name: name,
		Price: price,
		Description: &description,
		Inventory: inventory,
		Available: inventory > 0,
	}
	
	if err := s.db.WithContext(ctx).Create(product).Error; err != nil{
		return nil, err
	}
	return product, nil
}

func (s *ProductService)UpdateProduct(ctx context.Context, id string,  input models.UpdateProductInput) (*models.Product, error){
	product := &models.Product{ID: id}

	updates := s.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", id)

	if input.Name != nil{
		updates = updates.Update("name", *input.Name)
	}
	
	if input.Price != nil{
		updates = updates.Update("price", *input.Price)
	}
	if input.Description != nil{
		updates = updates.Update("description", *input.Description)
	}
	if input.Inventory != nil{
		updates = updates.Update("inventory", *input.Inventory)
	}

	// Execute the update
	if err := updates.Error; err != nil{
		return nil, err
	}

	// Return the updated user
	if err := s.db.WithContext(ctx).First(&product, "id = ?", id).Error; err != nil{
		return nil, err
	}
	return product, nil
}

func (s *ProductService)DeleteProduct(ctx context.Context, input models.DeleteProductInput) (bool, error){
	var result *gorm.DB

	if input.ID == nil && input.Name == nil {
		return false, errors.New("either ID or Name must be provided for deletion")
	}

	if input.ID != nil {
		result = s.db.WithContext(ctx).Delete(&models.Product{}, "id = ?", *input.ID)
	} else if input.Name != nil {
		result = s.db.WithContext(ctx).Delete(&models.Product{}, "name = ?", *input.Name)
	}

	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, fmt.Errorf("product not found")
	}
	return true, nil
}

func (s *ProductService)RestockProduct(ctx context.Context, id string, quantity int)(*models.Product, error) {
	var product models.Product

	// Fetching product
	if err := s.db.WithContext(ctx).First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Validate the restock amount
	if quantity <= 0 {
		return nil, fmt.Errorf("invalid restock amount: must be greater than zero")
	}
	// Update inventory
	product.Inventory += quantity

	// Save updated product
	if err := s.db.WithContext(ctx).Save(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil 
}

func (s *ProductService)SetProductAvailability(ctx context.Context, id string, available bool) (*models.Product, error){
	var product models.Product

	// Fetch product
	if err := s.db.WithContext(ctx).First(&product, "id = ?", id).Error; err != nil{
		return nil, err
	}

	// Prevent redundate updates
	if product.Available == available {
		return &product, fmt.Errorf("product already availability set to %t", available)

	}

	// safegard
	if available && product.Inventory <= 0{
		return nil, fmt.Errorf("cannot mark product as available with zero inventory")
	}

	// Udate availability
	product.Available = available

	if err := s.db.WithContext(ctx).Save(&product).Error; err != nil{
		return nil, err
	}
	return &product, nil
}

// GetAllProducts returns filderd products from db
func (s *ProductService) GetAllProductsCursor(ctx context.Context,  after *string, first int) ([]*models.Product, bool, error ){
	var products []*models.Product

    query := s.db.WithContext(ctx).
        Model(&models.Product{}).
        Order("id ASC")

    // If we have a cursor, decode it and filter
    if after != nil && *after != "" {
        lastID, err := DecodeCursor(*after)
        if err != nil {
            return nil, false, fmt.Errorf("invalid cursor: %w", err)
        }

        // Only return products *after* this ID
        query = query.Where("id > ?", lastID)
    }

    // Fetch first + 1 to check if there's a next page
    if err := query.Limit(first + 1).Find(&products).Error; err != nil {
        return nil, false, err
    }

    // Check if there are more results
    hasNextPage := len(products) > first
    if hasNextPage {
        products = products[:first]
    }

    return products, hasNextPage, nil
}
func (s *ProductService) CountProducts(ctx context.Context) (int, error) {
    var total int64
    err := s.db.Model(&models.Product{}).Count(&total).Error
    return int(total), err
}


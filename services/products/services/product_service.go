package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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
	product := &models.Product{
		ID:   fmt.Sprintf("product_%d", time.Now().UnixNano()),
		Name: name,
		Price: price,
		Description: &description,
		Inventory: inventory,
		Available: true,
	}
	if strings.TrimSpace(name) == ""{
		return nil, fmt.Errorf("invalid product name: missing or invalid field")
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

	// Update inventory
	product.Inventory += quantity

	if err := s.db.WithContext(ctx).Save(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil 
}

func (s *ProductService)SetProductAvailability(ctx context.Context, id string, available bool) (*models.Product, error){
	var product models.Product
	if err := s.db.WithContext(ctx).First(&product, "id = ?", id).Error; err != nil{
		return nil, err
	}

	product.Available = available

	if err := s.db.WithContext(ctx).Save(&product).Error; err != nil{
		return nil, err
	}
	return &product, nil
}
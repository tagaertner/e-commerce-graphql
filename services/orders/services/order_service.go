package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tagaertner/e-commerce-graphql/services/orders/models"
	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) GetAllOrders() ([]*models.Order, error) {
	var orders []*models.Order
	if err := s.db.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderService) GetOrderByID(id string) (*models.Order, error) {
	var order models.Order
	if err := s.db.First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) GetOrdersByUserID(userID string) ([]*models.Order, error) {
	var orders []*models.Order
	if err := s.db.
        Where("user_id = ?", userID).
        Preload("Products").
		Order("created_at ASC").
        Find(&orders).Error; err != nil {
        return nil, err
	}
	return orders, nil
}

func (s *OrderService)CreateOrder(ctx context.Context, userId string, productIds [] string, quantity int, totalPrice float64, status string, createdAt time.Time ) (*models.Order, error){
	order := &models.Order {
		ID: fmt.Sprintf("order_%d", time.Now().UnixNano()),
		UserID: userId,
		Quantity: quantity,
		TotalPrice: totalPrice,
		Status: status,
		CreatedAt: models.Time(createdAt),
	}

	if userId == "" || len(productIds) == 0 || quantity <= 0 || totalPrice <= 0 {
		return nil, errors.New("invalid order input: missing or invalid fields")
	}
	
	for _, pid := range productIds {
		order.Products = append(order.Products, models.Product{ID: pid})
	}
	if err := s.db.WithContext(ctx).Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, input *models.UpdateOrderInput) (*models.Order, error) {
	// Fetch existing order by ID
	var order models.Order
	if err := s.db.WithContext(ctx).First(&order, "id = ?", input.OrderID).Error; err != nil {
		return nil, err
	}

	// Apply updates only if the fields are not nil
	if input.Quantity != nil {
		order.Quantity = *input.Quantity
	}
	if input.TotalPrice != nil {
		order.TotalPrice = *input.TotalPrice
	}
	if input.Status != nil {
		order.Status = *input.Status
	}

	// Save the updated order
	if err := s.db.WithContext(ctx).Save(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *OrderService)DeleteOrder(ctx context.Context, input models.DeleteOrderInput) (bool, error) {
	var result *gorm.DB

	// Guard clause: require at least OrderID
	if input.OrderID == "" {
		return false, errors.New("OrderID must be provided for deletion")
	}

	// Delete by OrderID
	result = s.db.WithContext(ctx).Delete(&models.Order{}, "id = ?", input.OrderID)

	// Handle errors and no-op cases
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, errors.New("order not found")
	}
	return true, nil
}

func (s *OrderService)SetOrderStatus(ctx context.Context, input models.SetOrderStatusInput) (*models.Order, error) {
	order := &models.Order{ID: input.OrderID}

	if err := s.db.WithContext(ctx).
	Model(order).
	Update("status", input.Status).Error; err != nil {
		return nil, err
	}
	return order, nil
}

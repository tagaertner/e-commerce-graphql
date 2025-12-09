package services

import (
	"context"
	"fmt"
	"time"

	"github.com/tagaertner/e-commerce-graphql/services/users/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// Query
func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Mutation
func (s *UserService)CreateUser(ctx context.Context, name, email string, password string, role models.Role, active bool) (*models.User, error){
	user := &models.User{
		ID: fmt.Sprintf("user_%d", time.Now().UnixNano()),
		Name:   name,
		Email:  email,
		Password: password,
		Role:   role,  
		Active: true, 
	} 

	// TODO: Hash password before saving
	// hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err := s.db.WithContext(ctx).Create(user).Error; err !=nil {
		return nil, err
	}
	return user, nil

}

func (s *UserService) UpdateUser(ctx context.Context, input *models.UpdateUserInput) (*models.User, error) {
	var user models.User

	// Find the user first
	if err := s.db.WithContext(ctx).First(&user, "id = ?", input.ID).Error; err != nil {
		return nil, err
	}

	// Build update map with only provided fields
	updates := make(map[string]interface{})

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.Role != nil {
		updates["role"] = *input.Role
	}
	if input.Active != nil {
		updates["active"] = *input.Active
	}

	// Apply updates only if something to update
	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(&user).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s *UserService)DeleteUser(ctx context.Context, id string) (bool, error){
	result := s.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if result.Error != nil{
		return false, result.Error
	}

	return true, nil
}


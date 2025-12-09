package models

type Role string

const (
    RoleAdmin    Role = "ADMIN"
    RoleCustomer Role = "CUSTOMER"
)

type User struct {
    ID       string `json:"id" gorm:"primarykey"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Role     Role   `json:"role"`
    Active   bool   `json:"active"`
}

type CreateUserInput struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Role     Role   `json:"role"`
    Active   bool   `json:"active"`
}

type UpdateUserInput struct {
    ID     string  `json:"id"`
    Name   *string `json:"name"`
    Email  *string `json:"email"`
    Role   *Role   `json:"role"`
    Active *bool   `json:"active"`
}

func (User) IsEntity() {}
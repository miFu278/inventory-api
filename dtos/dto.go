package dtos

type TransactionType string

const (
	TransactionTypeIn  TransactionType = "IN"
	TransactionTypeOut TransactionType = "OUT"
)

// Product DTOs
type CreateProductInput struct {
	Name        string  `json:"name" minLength:"1" maxLength:"255" doc:"Product name"`
	SKU         string  `json:"sku" minLength:"1" maxLength:"100" doc:"Stock Keeping Unit"`
	Description string  `json:"description,omitempty" doc:"Product description"`
	Price       float64 `json:"price" minimum:"0.01" doc:"Product price (must be greater than 0)"`
	Quantity    int     `json:"quantity" minimum:"1" doc:"Initial quantity (must be at least 1)"`
}

type UpdateProductInput struct {
	Name        *string  `json:"name,omitempty" minLength:"1" maxLength:"255" doc:"Product name"`
	Description *string  `json:"description,omitempty" doc:"Product description"`
	Price       *float64 `json:"price,omitempty" minimum:"0" doc:"Product price (can be 0)"`
	Quantity    *int     `json:"quantity,omitempty" minimum:"0" doc:"Product quantity (can be 0)"`
}

type ProductResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	SKU         string  `json:"sku"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Transaction DTOs
type CreateTransactionInput struct {
	ProductID       uint            `json:"product_id" doc:"Product ID"`
	Quantity        int             `json:"quantity" minimum:"1" doc:"Transaction quantity"`
	TransactionType TransactionType `json:"transaction_type" enum:"IN,OUT" doc:"Transaction type (IN/OUT)"`
	Notes           string          `json:"notes,omitempty" doc:"Transaction notes"`
}

type TransactionResponse struct {
	ID              uint             `json:"id"`
	ProductID       uint             `json:"product_id"`
	Product         *ProductResponse `json:"product,omitempty"`
	Quantity        int              `json:"quantity"`
	TransactionType TransactionType  `json:"transaction_type"`
	Notes           string           `json:"notes"`
	CreatedAt       string           `json:"created_at"`
	UpdatedAt       string           `json:"updated_at"`
}

type ProductFilter struct {
	SKU      *string  `json:"sku,omitempty"`
	Name     *string  `json:"name,omitempty"`
	MinPrice *float64 `json:"min_price,omitempty"`
	MaxPrice *float64 `json:"max_price,omitempty"`
}

func (f *ProductFilter) IsEmpty() bool {
	if f == nil {
		return true
	}
	return f.SKU == nil && f.Name == nil && f.MinPrice == nil && f.MaxPrice == nil
}

func (f *ProductFilter) HasSKU() bool {
	return f != nil && f.SKU != nil && *f.SKU != ""
}

func (f *ProductFilter) HasName() bool {
	return f != nil && f.Name != nil && *f.Name != ""
}

func (f *ProductFilter) HasMinPrice() bool {
	return f != nil && f.MinPrice != nil
}

func (f *ProductFilter) HasMaxPrice() bool {
	return f != nil && f.MaxPrice != nil
}

// User DTOs
type RegisterInput struct {
	Username string `json:"username" minLength:"3" maxLength:"50" pattern:"^[a-zA-Z0-9_]+$" doc:"Username (alphanumeric and underscore only)"`
	Password string `json:"password" minLength:"8" maxLength:"100" doc:"Password (minimum 8 characters)"`
	Email    string `json:"email" format:"email" doc:"Email address"`
	Phone    string `json:"phone" minLength:"10" maxLength:"15" pattern:"^[0-9+\\-\\s()]+$" doc:"Phone number"`
	Role     string `json:"role" enum:"admin,user" default:"user" doc:"User role (admin/user)"`
}

type LoginInput struct {
	Username string `json:"username" minLength:"1" doc:"Username"`
	Password string `json:"password" minLength:"1" doc:"Password"`
}

type UpdateUserInput struct {
	Email *string `json:"email,omitempty" format:"email" doc:"Email address"`
	Phone *string `json:"phone,omitempty" minLength:"10" maxLength:"15" doc:"Phone number"`
	Role  *string `json:"role,omitempty" enum:"admin,user" doc:"User role (admin/user)"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" minLength:"1" doc:"Current password"`
	NewPassword string `json:"new_password" minLength:"8" maxLength:"100" doc:"New password (minimum 8 characters)"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

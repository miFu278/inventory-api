package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:255" json:"name"`
	SKU         string         `gorm:"uniqueIndex;not null;size:100" json:"sku"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Quantity    int            `gorm:"not null;default:0" json:"quantity"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type TransactionType string

const (
	TransactionTypeIn  TransactionType = "IN"
	TransactionTypeOut TransactionType = "OUT"
)

type Transaction struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	ProductID       uint            `gorm:"not null;index" json:"product_id"`
	Product         Product         `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"product,omitempty"`
	Quantity        int             `gorm:"not null" json:"quantity"`
	TransactionType TransactionType `gorm:"not null;size:10" json:"transaction_type"`
	Notes           string          `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// DTOs for API
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

type CreateTransactionInput struct {
	ProductID       uint            `json:"product_id" doc:"Product ID"`
	Quantity        int             `json:"quantity" minimum:"1" doc:"Transaction quantity"`
	TransactionType TransactionType `json:"transaction_type" enum:"IN,OUT" doc:"Transaction type (IN/OUT)"`
	Notes           string          `json:"notes,omitempty" doc:"Transaction notes"`
}

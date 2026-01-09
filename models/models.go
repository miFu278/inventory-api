package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeIn  TransactionType = "IN"
	TransactionTypeOut TransactionType = "OUT"
)

type Product struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"not null;size:255"`
	SKU         string  `gorm:"uniqueIndex;not null;size:100"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"type:decimal(10,2);not null"`
	Quantity    int     `gorm:"not null;default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Transaction struct {
	ID              uint            `gorm:"primaryKey"`
	ProductID       uint            `gorm:"not null;index"`
	Product         Product         `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Quantity        int             `gorm:"not null"`
	TransactionType TransactionType `gorm:"not null;size:10"`
	Notes           string          `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

package dtos

import (
	"inventory-api/models"
)

// ToProductModel converts CreateProductInput to Product model
func (dto *CreateProductInput) ToProductModel() *models.Product {
	return &models.Product{
		Name:        dto.Name,
		SKU:         dto.SKU,
		Description: dto.Description,
		Price:       dto.Price,
		Quantity:    dto.Quantity,
	}
}

// ToProductResponse converts Product model to ProductResponse DTO
func ToProductResponse(product *models.Product) *ProductResponse {
	if product == nil {
		return nil
	}
	return &ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		SKU:         product.SKU,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToProductResponseList converts slice of Product models to slice of ProductResponse DTOs
func ToProductResponseList(products []models.Product) []ProductResponse {
	responses := make([]ProductResponse, len(products))
	for i, product := range products {
		responses[i] = *ToProductResponse(&product)
	}
	return responses
}

// ApplyToProduct applies UpdateProductInput to existing Product model
func (dto *UpdateProductInput) ApplyToProduct(product *models.Product) {
	if dto.Name != nil {
		product.Name = *dto.Name
	}
	if dto.Description != nil {
		product.Description = *dto.Description
	}
	if dto.Price != nil {
		product.Price = *dto.Price
	}
}

// ToTransactionModel converts CreateTransactionInput to Transaction model
func (dto *CreateTransactionInput) ToTransactionModel() *models.Transaction {
	return &models.Transaction{
		ProductID:       dto.ProductID,
		Quantity:        dto.Quantity,
		TransactionType: models.TransactionType(dto.TransactionType),
		Notes:           dto.Notes,
	}
}

// ToTransactionResponse converts Transaction model to TransactionResponse DTO
func ToTransactionResponse(transaction *models.Transaction) *TransactionResponse {
	if transaction == nil {
		return nil
	}

	response := &TransactionResponse{
		ID:              transaction.ID,
		ProductID:       transaction.ProductID,
		Quantity:        transaction.Quantity,
		TransactionType: TransactionType(transaction.TransactionType),
		Notes:           transaction.Notes,
		CreatedAt:       transaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       transaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Include product if loaded
	if transaction.Product.ID != 0 {
		response.Product = ToProductResponse(&transaction.Product)
	}

	return response
}

// ToTransactionResponseList converts slice of Transaction models to slice of TransactionResponse DTOs
func ToTransactionResponseList(transactions []models.Transaction) []TransactionResponse {
	responses := make([]TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = *ToTransactionResponse(&transaction)
	}
	return responses
}

// ToUserResponse converts User model to UserResponse DTO
func ToUserResponse(user *models.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
	}
}

// ToUserResponseList converts slice of User models to slice of UserResponse DTOs
func ToUserResponseList(users []models.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = *ToUserResponse(&user)
	}
	return responses
}

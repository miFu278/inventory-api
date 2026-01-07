package services

import (
	"errors"
	"inventory-api/internal/models"
	"inventory-api/internal/repo"

	"gorm.io/gorm"
)

type InventoryService struct {
	repo *repo.InventoryRepository
}

func NewInventoryService(repo *repo.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

// Product services
func (s *InventoryService) CreateProduct(input *models.CreateProductInput) (*models.Product, error) {
	// Check if SKU already exists
	existing, err := s.repo.GetProductBySKU(input.SKU)
	if err == nil && existing != nil {
		return nil, errors.New("product with this SKU already exists")
	}

	product := &models.Product{
		Name:        input.Name,
		SKU:         input.SKU,
		Description: input.Description,
		Price:       input.Price,
		Quantity:    input.Quantity,
	}

	if err := s.repo.CreateProduct(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *InventoryService) GetProductByID(id uint) (*models.Product, error) {
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return product, nil
}

func (s *InventoryService) GetProductBySKU(sku string) (*models.Product, error) {
	product, err := s.repo.GetProductBySKU(sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return product, nil
}

func (s *InventoryService) GetAllProducts(limit, offset int) ([]models.Product, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetAllProducts(limit, offset)
}

func (s *InventoryService) UpdateProduct(id uint, input *models.UpdateProductInput) (*models.Product, error) {
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Update only provided fields
	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Price != nil {
		product.Price = *input.Price
	}

	if err := s.repo.UpdateProduct(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *InventoryService) DeleteProduct(id uint) error {
	// Check if product exists
	_, err := s.repo.GetProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	return s.repo.DeleteProduct(id)
}

// Transaction services
func (s *InventoryService) CreateTransaction(input *models.CreateTransactionInput) (*models.Transaction, error) {
	// Validate product exists
	_, err := s.repo.GetProductByID(input.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Validate quantity
	if input.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	// Update product quantity with transaction
	err = s.repo.UpdateProductQuantityWithTransaction(
		input.ProductID,
		input.Quantity,
		input.TransactionType,
		input.Notes,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			return nil, errors.New("insufficient quantity for OUT transaction")
		}
		return nil, err
	}

	// Get the created transaction (last one for this product)
	transactions, err := s.repo.GetTransactionsByProductID(input.ProductID, 1, 0)
	if err != nil || len(transactions) == 0 {
		return nil, errors.New("failed to retrieve created transaction")
	}

	return &transactions[0], nil
}

func (s *InventoryService) GetTransactionByID(id uint) (*models.Transaction, error) {
	transaction, err := s.repo.GetTransactionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return transaction, nil
}

func (s *InventoryService) GetTransactionsByProductID(productID uint, limit, offset int) ([]models.Transaction, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Validate product exists
	_, err := s.repo.GetProductByID(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return s.repo.GetTransactionsByProductID(productID, limit, offset)
}

func (s *InventoryService) GetAllTransactions(limit, offset int) ([]models.Transaction, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetAllTransactions(limit, offset)
}

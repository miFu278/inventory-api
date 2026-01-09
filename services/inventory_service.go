package services

import (
	"errors"
	"inventory-api/dtos"
	"inventory-api/models"
	"inventory-api/repo"

	"gorm.io/gorm"
)

type InventoryService struct {
	repo *repo.InventoryRepository
}

func NewInventoryService(repo *repo.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

// Product services
func (s *InventoryService) CreateProduct(input *dtos.CreateProductInput) (*dtos.ProductResponse, error) {
	// Check if SKU already exists
	existing, err := s.repo.GetProductBySKU(input.SKU)
	if err == nil && existing != nil {
		return nil, errors.New("product with this SKU already exists")
	}

	// Convert DTO to model
	product := input.ToProductModel()

	if err := s.repo.CreateProduct(product); err != nil {
		return nil, err
	}

	// Convert model to response DTO
	return dtos.ToProductResponse(product), nil
}

func (s *InventoryService) GetProductByID(id uint) (*dtos.ProductResponse, error) {
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return dtos.ToProductResponse(product), nil
}

func (s *InventoryService) GetProductBySKU(sku string) (*dtos.ProductResponse, error) {
	product, err := s.repo.GetProductBySKU(sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return dtos.ToProductResponse(product), nil
}

func (s *InventoryService) GetAllProducts(limit, offset int) ([]dtos.ProductResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	products, err := s.repo.GetAllProducts(limit, offset)
	if err != nil {
		return nil, err
	}

	return dtos.ToProductResponseList(products), nil
}

// GetProductsWithFilter retrieves products with filtering
func (s *InventoryService) GetProductsWithFilter(filter *dtos.ProductFilter, limit, offset int) ([]dtos.ProductResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	products, err := s.repo.GetProductsWithFilter(filter, limit, offset)
	if err != nil {
		return nil, err
	}

	return dtos.ToProductResponseList(products), nil
}

func (s *InventoryService) UpdateProduct(id uint, input *dtos.UpdateProductInput) (*dtos.ProductResponse, error) {
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Apply DTO updates to model
	input.ApplyToProduct(product)

	if err := s.repo.UpdateProduct(product); err != nil {
		return nil, err
	}

	return dtos.ToProductResponse(product), nil
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
func (s *InventoryService) CreateTransaction(input *dtos.CreateTransactionInput) (*dtos.TransactionResponse, error) {
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
		models.TransactionType(input.TransactionType),
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

	return dtos.ToTransactionResponse(&transactions[0]), nil
}

func (s *InventoryService) GetTransactionByID(id uint) (*dtos.TransactionResponse, error) {
	transaction, err := s.repo.GetTransactionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return dtos.ToTransactionResponse(transaction), nil
}

func (s *InventoryService) GetTransactionsByProductID(productID uint, limit, offset int) ([]dtos.TransactionResponse, error) {
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

	transactions, err := s.repo.GetTransactionsByProductID(productID, limit, offset)
	if err != nil {
		return nil, err
	}

	return dtos.ToTransactionResponseList(transactions), nil
}

func (s *InventoryService) GetAllTransactions(limit, offset int) ([]dtos.TransactionResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	transactions, err := s.repo.GetAllTransactions(limit, offset)
	if err != nil {
		return nil, err
	}

	return dtos.ToTransactionResponseList(transactions), nil
}

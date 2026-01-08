package repo

import (
	"inventory-api/models"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// Product operations
func (r *InventoryRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *InventoryRepository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *InventoryRepository) GetProductBySKU(sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *InventoryRepository) GetAllProducts(limit, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Limit(limit).Offset(offset).Find(&products).Error
	return products, err
}

func (r *InventoryRepository) UpdateProduct(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *InventoryRepository) DeleteProduct(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

// Transaction operations
func (r *InventoryRepository) CreateTransaction(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *InventoryRepository) GetTransactionByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Product").First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *InventoryRepository) GetTransactionsByProductID(productID uint, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("product_id = ?", productID).
		Preload("Product").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *InventoryRepository) GetAllTransactions(limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Product").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

// Update product quantity with transaction
func (r *InventoryRepository) UpdateProductQuantityWithTransaction(productID uint, quantity int, txType models.TransactionType, notes string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get product
		var product models.Product
		if err := tx.First(&product, productID).Error; err != nil {
			return err
		}

		// Update quantity
		if txType == models.TransactionTypeIn {
			product.Quantity += quantity
		} else {
			if product.Quantity < quantity {
				return gorm.ErrInvalidData
			}
			product.Quantity -= quantity
		}

		// Save product
		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		// Create transaction record
		transaction := models.Transaction{
			ProductID:       productID,
			Quantity:        quantity,
			TransactionType: txType,
			Notes:           notes,
		}
		return tx.Create(&transaction).Error
	})
}

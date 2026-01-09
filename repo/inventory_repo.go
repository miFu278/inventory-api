package repo

import (
	"context"
	"inventory-api/dtos"
	"inventory-api/models"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	db              *gorm.DB
	productRepo     *BaseRepository[models.Product]
	transactionRepo *BaseRepository[models.Transaction]
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{
		db:              db,
		productRepo:     NewBaseRepository[models.Product](db),
		transactionRepo: NewBaseRepository[models.Transaction](db),
	}
}

// Product operations using BaseRepository
func (r *InventoryRepository) CreateProduct(product *models.Product) error {
	return r.productRepo.Create(context.Background(), product)
}

func (r *InventoryRepository) GetProductByID(id uint) (*models.Product, error) {
	return r.productRepo.GetByID(context.Background(), id)
}

func (r *InventoryRepository) GetProductBySKU(sku string) (*models.Product, error) {
	return r.productRepo.FindOne(context.Background(), func(db *gorm.DB) *gorm.DB {
		return db.Where("sku = ?", sku)
	})
}

func (r *InventoryRepository) GetAllProducts(limit, offset int) ([]models.Product, error) {
	return r.productRepo.List(
		context.Background(),
		WithLimit(limit),
		WithOffset(offset),
	)
}

// GetProductsWithFilter retrieves products with filtering support
func (r *InventoryRepository) GetProductsWithFilter(filter *dtos.ProductFilter, limit, offset int) ([]models.Product, error) {
	scopes := r.buildProductFilterScopes(filter)

	// Add pagination
	scopes = append(scopes, WithLimit(limit), WithOffset(offset))

	return r.productRepo.List(context.Background(), scopes...)
}

// buildProductFilterScopes converts ProductFilter to GORM scopes
func (r *InventoryRepository) buildProductFilterScopes(filter *dtos.ProductFilter) []func(*gorm.DB) *gorm.DB {
	scopes := []func(*gorm.DB) *gorm.DB{}

	if filter == nil || filter.IsEmpty() {
		return scopes
	}

	// Filter by SKU (exact match)
	if filter.HasSKU() {
		sku := *filter.SKU
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("sku = ?", sku)
		})
	}

	// Filter by Name (partial match, case-insensitive)
	if filter.HasName() {
		name := *filter.Name
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("name ILIKE ?", "%"+name+"%")
		})
	}

	// Filter by MinPrice
	if filter.HasMinPrice() {
		minPrice := *filter.MinPrice
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("price >= ?", minPrice)
		})
	}

	// Filter by MaxPrice
	if filter.HasMaxPrice() {
		maxPrice := *filter.MaxPrice
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("price <= ?", maxPrice)
		})
	}

	return scopes
}

func (r *InventoryRepository) UpdateProduct(product *models.Product) error {
	return r.productRepo.Update(context.Background(), product)
}

func (r *InventoryRepository) DeleteProduct(id uint) error {
	return r.productRepo.Delete(context.Background(), id)
}

// Transaction operations using BaseRepository
func (r *InventoryRepository) CreateTransaction(tx *models.Transaction) error {
	return r.transactionRepo.Create(context.Background(), tx)
}

func (r *InventoryRepository) GetTransactionByID(id uint) (*models.Transaction, error) {
	return r.transactionRepo.FindOne(
		context.Background(),
		func(db *gorm.DB) *gorm.DB {
			return db.Where("id = ?", id)
		},
		WithPreload("Product"),
	)
}

func (r *InventoryRepository) GetTransactionsByProductID(productID uint, limit, offset int) ([]models.Transaction, error) {
	return r.transactionRepo.List(
		context.Background(),
		func(db *gorm.DB) *gorm.DB {
			return db.Where("product_id = ?", productID)
		},
		WithPreload("Product"),
		WithLimit(limit),
		WithOffset(offset),
		WithOrder("created_at DESC"),
	)
}

func (r *InventoryRepository) GetAllTransactions(limit, offset int) ([]models.Transaction, error) {
	return r.transactionRepo.List(
		context.Background(),
		WithPreload("Product"),
		WithLimit(limit),
		WithOffset(offset),
		WithOrder("created_at DESC"),
	)
}

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

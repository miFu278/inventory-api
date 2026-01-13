package repo

import "inventory-api/models"

// UserRepositoryInterface định nghĩa contract cho UserRepository
// Giúp dễ dàng mock trong tests
type UserRepositoryInterface interface {
	GetUserByID(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByPhone(phone string) (*models.User, error)
	GetAllUsers(limit, offset int) ([]models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id uint) error
}

// Verify UserRepository implements the interface
var _ UserRepositoryInterface = (*UserRepository)(nil)

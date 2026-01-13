package repo

import (
	"context"
	"inventory-api/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db       *gorm.DB
	userRepo *BaseRepository[models.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db:       db,
		userRepo: NewBaseRepository[models.User](db),
	}
}

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	return r.userRepo.GetByID(context.Background(), id)
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	return r.userRepo.FindOne(context.Background(), func(db *gorm.DB) *gorm.DB {
		return db.Where("username = ?", username)
	})
}

func (r *UserRepository) GetAllUsers(limit, offset int) ([]models.User, error) {
	return r.userRepo.List(
		context.Background(),
		WithLimit(limit),
		WithOffset(offset),
	)
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	return r.userRepo.FindOne(context.Background(), func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", email)
	})
}

func (r *UserRepository) GetUserByPhone(phone string) (*models.User, error) {
	return r.userRepo.FindOne(context.Background(), func(db *gorm.DB) *gorm.DB {
		return db.Where("phone = ?", phone)
	})
}

func (r *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	err := r.userRepo.Create(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user *models.User) (*models.User, error) {
	err := r.userRepo.Update(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) DeleteUser(id uint) error {
	return r.userRepo.Delete(context.Background(), id)
}

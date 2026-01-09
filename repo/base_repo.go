package repo

import (
	"context"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

func (r *BaseRepository[T]) List(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) ([]T, error) {
	var results []T
	err := r.db.WithContext(ctx).Scopes(scopes...).Find(&results).Error
	return results, err
}

func (r *BaseRepository[T]) Count(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var count int64
	var entity T
	err := r.db.WithContext(ctx).Model(&entity).Scopes(scopes...).Count(&count).Error
	return count, err
}

func (r *BaseRepository[T]) FindOne(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).Scopes(scopes...).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func WithLimit(limit int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func WithOffset(offset int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

func WithOrder(order string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(order)
	}
}

func WithPreload(associations ...string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, assoc := range associations {
			db = db.Preload(assoc)
		}
		return db
	}
}

// WithWhere adds a WHERE condition
func WithWhere(query interface{}, args ...interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	}
}

package order

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, order *Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) Repository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(ctx context.Context, order *Order) error {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		return fmt.Errorf("failed to create order: %v", err)
	}

	return nil
}

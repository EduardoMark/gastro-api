package dishes

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dish *Dish) error
	GetOneByID(ctx context.Context, id uuid.UUID) (*Dish, error)
	GetOneByName(ctx context.Context, name string) (*Dish, error)
	Query(ctx context.Context) ([]*Dish, error)
	Update(ctx context.Context, dish *Dish) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type dishRepository struct {
	db *gorm.DB
}

func NewDishRepository(db *gorm.DB) Repository {
	return &dishRepository{
		db: db,
	}
}

var ErrDishAlreadyExists = errors.New("dish already exists")
var ErrDishNotFound = errors.New("dish not found")

func (r *dishRepository) Create(ctx context.Context, dish *Dish) error {
	err := r.db.WithContext(ctx).Create(dish).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDishAlreadyExists
		}
		return fmt.Errorf("Create - failed to create user: %v", err)
	}
	return nil
}

func (r *dishRepository) GetOneByID(ctx context.Context, id uuid.UUID) (*Dish, error) {
	var dish Dish

	err := r.db.WithContext(ctx).First(&dish, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDishNotFound
		}
		return nil, fmt.Errorf("GetOneByID - failed to get dish: %v", err)
	}

	return &dish, nil
}

func (r *dishRepository) GetOneByName(ctx context.Context, name string) (*Dish, error) {
	var dish Dish

	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&dish).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDishNotFound
		}
		return nil, fmt.Errorf("GetOneByName - failed to get dish: %v", err)
	}

	return &dish, nil
}

func (r *dishRepository) Query(ctx context.Context) ([]*Dish, error) {
	var dishes []*Dish

	err := r.db.WithContext(ctx).Find(&dishes).Error
	if err != nil {
		return nil, fmt.Errorf("Query - failed find all dishes: %v", err)
	}

	if len(dishes) == 0 {
		return nil, ErrDishNotFound
	}

	return dishes, nil
}

func (r *dishRepository) Update(ctx context.Context, dish *Dish) error {
	result := r.db.WithContext(ctx).Model(Dish{}).Where("id = ?", dish.ID).Updates(dish)
	if result.Error != nil {
		return fmt.Errorf("Update - failed to update dish: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrDishNotFound
	}

	return nil
}

func (r *dishRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(Dish{}).Where("id = ?", id).Delete(nil)

	if result.Error != nil {
		return fmt.Errorf("Delete - failed to delete dish: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrDishNotFound
	}

	return nil
}

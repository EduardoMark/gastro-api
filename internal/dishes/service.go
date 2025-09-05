package dishes

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Service interface {
	Create(ctx context.Context, name, description, category string, price float64) error
	GetOneByID(ctx context.Context, id uuid.UUID) (*Dish, error)
	Query(ctx context.Context) ([]*Dish, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type dishService struct {
	r Repository
}

func NewDishService(r Repository) Service {
	return &dishService{
		r: r,
	}
}

func (s *dishService) Create(ctx context.Context, name, description, category string, price float64) error {
	decimalPrice, err := decimal.NewFromString(fmt.Sprintf("%.2f", price))
	if err != nil {
		return fmt.Errorf("failed to convert price to decimal price: %v", err)
	}

	dish := Dish{
		Name:        name,
		Description: description,
		Category:    category,
		Price:       decimalPrice,
	}

	if err := s.r.Create(ctx, &dish); err != nil {
		return err
	}

	return nil
}

func (s *dishService) GetOneByID(ctx context.Context, id uuid.UUID) (*Dish, error) {
	record, err := s.r.GetOneByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (s *dishService) Query(ctx context.Context) ([]*Dish, error) {
	records, err := s.r.Query(ctx)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *dishService) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) error {
	decimalPrice, err := decimal.NewFromString(fmt.Sprintf("%.2f", req.Price))
	if err != nil {
		return fmt.Errorf("failed to convert price to decimal price: %v", err)
	}

	dish := Dish{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       decimalPrice,
		Category:    req.Category,
	}

	if err := s.r.Update(ctx, &dish); err != nil {
		return err
	}

	return nil
}

func (s *dishService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.r.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

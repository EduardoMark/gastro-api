package order

import (
	"context"
	"fmt"

	"github.com/EduardoMark/gastro-api/internal/dishes"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, items []createOrderItems) error
}

type orderService struct {
	repository Repository
	dishRepo   dishes.Repository
}

func NewOrderService(repository Repository, dishRepo dishes.Repository) Service {
	return &orderService{
		repository: repository,
		dishRepo:   dishRepo,
	}
}

func (s *orderService) Create(ctx context.Context, userID uuid.UUID, items []createOrderItems) error {
	order := Order{
		UserID:      userID,
		Status:      STATUS_NEW,
		Items:       []OrderItem{},
		TotalAmount: decimal.NewFromInt(0),
	}

	for _, i := range items {
		dishID, err := uuid.Parse(i.DishID)
		if err != nil {
			return fmt.Errorf("error on parse dish id to uuid type: %v", err)
		}

		dish, err := s.dishRepo.GetOneByID(ctx, dishID)
		if err != nil {
			return err
		}

		price := dish.Price
		subTotal := price.Mul(decimal.NewFromInt(int64(i.Quantity)))

		order.Items = append(order.Items, OrderItem{
			DishID:   dishID,
			Quantity: i.Quantity,
			Price:    price,
			SubTotal: subTotal,
		})

		order.TotalAmount = order.TotalAmount.Add(subTotal)
	}

	if err := s.repository.Create(ctx, &order); err != nil {
		return err
	}

	return nil
}

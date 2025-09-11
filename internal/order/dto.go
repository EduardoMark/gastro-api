package order

import (
	"fmt"

	"github.com/EduardoMark/gastro-api/internal/validation"
	"github.com/go-playground/validator/v10"
)

type CreateOrderRequest struct {
	Items []createOrderItems `json:"items" validete:"required,min=1"`
}

type createOrderItems struct {
	DishID   string `json:"dish_id"`
	Quantity int    `json:"quantity"`
}

func (r *CreateOrderRequest) Validate() error {
	if err := validation.Validate.Struct(r); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Tag() == "required" {
				return fmt.Errorf("field %s is required", err.Field())
			}
			if err.Tag() == "min" {
				return fmt.Errorf("field %s must be at least %s item", err.Field(), err.Param())
			}
		}
	}
	return nil
}

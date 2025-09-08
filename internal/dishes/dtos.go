package dishes

import (
	"fmt"
	"time"

	"github.com/EduardoMark/gastro-api/internal/validation"
	"github.com/go-playground/validator/v10"
)

type CreateRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description string  `json:"description" validate:"required,min=3,max=500"`
	Price       float64 `json:"price" validate:"gt=0"`
	Category    string  `json:"category" validate:"required,min=3,max=100"`
}

func (r *CreateRequest) Validate() error {
	if err := validation.Validate.Struct(r); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Tag() == "required" {
				return fmt.Errorf("field %s is required", err.Field())
			}
			if err.Tag() == "min" {
				return fmt.Errorf("field %s must be at least %s characters long", err.Field(), err.Param())
			}
			if err.Tag() == "max" {
				return fmt.Errorf("field %s must be at most %s characters long", err.Field(), err.Param())
			}
			if err.Tag() == "gt" {
				return fmt.Errorf("field %s need a greater than 0 value", err.Field())
			}
		}
	}

	return nil
}

type DishResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description string  `json:"description" validate:"required,min=3,max=500"`
	Price       float64 `json:"price" validate:"gt=0"`
	Category    string  `json:"category" validate:"required,min=3,max=100"`
}

func (r *UpdateRequest) Validate() error {
	if err := validation.Validate.Struct(r); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Tag() == "required" {
				return fmt.Errorf("field %s is required", err.Field())
			}
			if err.Tag() == "min" {
				return fmt.Errorf("field %s must be at least %s characters long", err.Field(), err.Param())
			}
			if err.Tag() == "max" {
				return fmt.Errorf("field %s must be at most %s characters long", err.Field(), err.Param())
			}
			if err.Tag() == "gt" {
				return fmt.Errorf("field %s need a greater than 0 value", err.Field())
			}
		}
	}

	return nil
}

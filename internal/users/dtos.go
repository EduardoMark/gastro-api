package users

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type SignupRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

func (r SignupRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
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
			if err.Tag() == "email" {
				return fmt.Errorf("field %s must be a valid email address", err.Field())
			}
		}
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

func (r LoginRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Tag() == "required" {
				return fmt.Errorf("field %s is required", err.Field())
			}
			if err.Tag() == "email" {
				return fmt.Errorf("field %s must be a valid email address", err.Field())
			}
			if err.Tag() == "min" {
				return fmt.Errorf("field %s must be at least %s characters long", err.Field(), err.Param())
			}
			if err.Tag() == "max" {
				return fmt.Errorf("field %s must be at most %s characters long", err.Field(), err.Param())
			}
		}
	}
	return nil
}

type ChangePasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}

func (r ChangePasswordRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
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
		}
	}
	return nil
}

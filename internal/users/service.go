package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Authenticate(ctx context.Context, email, password string) (*User, error)
	Create(ctx context.Context, name, email, password string, role Role) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, newPassoword string) error
}

type userService struct {
	r Repository
}

func NewUserService(r Repository) Service {
	return &userService{
		r: r,
	}
}

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrSamePassword = errors.New("new password cannot be the same as the old password")

func (s *userService) Authenticate(ctx context.Context, email, password string) (*User, error) {
	user, err := s.r.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *userService) Create(ctx context.Context, name, email, password string, role Role) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user := User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	if role == RoleAdmin {
		user.Role = RoleAdmin
	}

	if role == RoleClient {
		user.Role = RoleClient
	}

	if err := s.r.CreateUser(ctx, &user); err != nil {
		return err
	}

	return nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.r.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

func (s *userService) ChangePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	user, err := s.r.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(newPassword)) == nil {
		return ErrSamePassword
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	err = s.r.UpdatePassword(ctx, userID, string(passwordHash))
	if err != nil {
		return err
	}

	return nil
}

package users

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// Deve implementar os mesmos metodos do repository real
type MockRepository struct {
	createUserFunc     func(ctx context.Context, user *User) error
	getUserByIDFunc    func(ctx context.Context, id uuid.UUID) (*User, error)
	getUserByEmailFunc func(ctx context.Context, email string) (*User, error)
	updatePasswordFunc func(ctx context.Context, id uuid.UUID, newHash string) error
	deleteUserFunc     func(ctx context.Context, id uuid.UUID) error
}

func (m *MockRepository) CreateUser(ctx context.Context, user *User) error {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, user)
	}
	return nil
}

func (m *MockRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockRepository) UpdatePassword(ctx context.Context, id uuid.UUID, newHash string) error {
	if m.updatePasswordFunc != nil {
		return m.updatePasswordFunc(ctx, id, newHash)
	}
	return nil
}

func (m *MockRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if m.deleteUserFunc != nil {
		return m.deleteUserFunc(ctx, id)
	}
	return nil
}

// TESTS

func TestCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("should create user successfully", func(t *testing.T) {
		// create mock
		mockRepo := &MockRepository{
			createUserFunc: func(ctx context.Context, user *User) error {
				return nil
			},
		}

		s := NewUserService(mockRepo)

		err := s.Create(ctx, "Eduardo", "eduardo@email.com", "12345678", RoleAdmin)
		if err != nil {
			t.Errorf("expected no error, but got: %v", err)
		}
	})

	t.Run("should return ErrEmailAlreadyExists when repo fails", func(t *testing.T) {
		mockRepo := &MockRepository{
			createUserFunc: func(ctx context.Context, user *User) error {
				return ErrEmailAlreadyExists
			},
		}

		s := NewUserService(mockRepo)

		err := s.Create(ctx, "Eduardo", "eduardo@email.com", "12345678", RoleClient)
		if !errors.Is(err, ErrEmailAlreadyExists) {
			t.Errorf("expected ErrEmailAlreadyExists, got: %v", err)
		}
	})

	t.Run("should return wrapped DB error when repo fails", func(t *testing.T) {
		mockRepo := &MockRepository{
			createUserFunc: func(ctx context.Context, user *User) error {
				return fmt.Errorf("failed to create user")
			},
		}

		s := NewUserService(mockRepo)

		err := s.Create(ctx, "Eduardo", "eduardo@email.com", "12345678", RoleClient)
		if err == nil {
			t.Errorf("expected Generic Database Error, got: %v", err)
		}

		expectedMsg := "failed to create user"
		if !strings.Contains(err.Error(), expectedMsg) {
			t.Errorf("expected error to contain %q, got: %v", expectedMsg, err)
		}
	})
}

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()

	t.Run("success on get user by email", func(t *testing.T) {
		mockRepo := &MockRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (*User, error) {
				return &User{
					ID:    uuid.UUID{},
					Name:  "Eduardo",
					Email: "eduardo@email.com",
					Role:  RoleAdmin,
				}, nil
			},
		}

		s := NewUserService(mockRepo)

		user, err := s.GetUserByEmail(ctx, "eduardo@email.com")
		if err != nil {
			t.Errorf("expected no error, but got: %v", err)
		}

		if user.Name != "Eduardo" || user.Email != "eduardo@email.com" || user.Role != RoleAdmin {
			t.Errorf("invalid user fields")
		}
	})
}

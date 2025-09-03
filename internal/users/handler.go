package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/EduardoMark/gastro-api/internal/auth"
	"github.com/EduardoMark/gastro-api/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	s             Service
	jwtMiddleware *middleware.JWTMiddleware
	authService   *auth.AuthJWTService
}

func NerUserHandler(
	s Service,
	jwtMiddleware *middleware.JWTMiddleware,
	authService *auth.AuthJWTService,
) UserHandler {
	return UserHandler{
		s:             s,
		jwtMiddleware: jwtMiddleware,
		authService:   authService,
	}
}

func (h *UserHandler) UserRoutes(r chi.Router) {
	r.Post("/login", h.Login)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.Signup)

		r.Group(func(r chi.Router) {
			r.Use(h.jwtMiddleware.JWTAuth)

			r.Put("/change-password", h.ChangePassword)
		})
	})
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Signup handler running...")
	ctx := r.Context()

	var body SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if err := body.Validate(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := h.s.Create(
		ctx,
		body.Name,
		body.Email,
		body.Password,
	); err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "email already exists",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"success": "user created with success",
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Login handler running...")
	ctx := r.Context()

	var body LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid body request",
		})
		return
	}

	if err := body.Validate(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	user, err := h.s.Authenticate(ctx, body.Email, body.Password)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "user not found",
			})
			return
		}

		if errors.Is(err, ErrInvalidCredentials) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid credentials",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	token, err := h.authService.New(user.ID.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Change Password handler running...")
	ctx := r.Context()

	userIDRaw := ctx.Value(middleware.CtxUserId)
	userID, ok := userIDRaw.(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid user id in context",
		})
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid user id type uuuid",
		})
		return
	}

	var body ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid body request",
		})
		return
	}

	if err := body.Validate(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := h.s.ChangePassword(ctx, uid, body.NewPassword); err != nil {
		if errors.Is(err, ErrSamePassword) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "new password cannot be the same as the old password",
			})
			return
		}

		if errors.Is(err, ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "user not found",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

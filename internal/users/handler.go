package users

import (
	"errors"
	"net/http"

	"github.com/EduardoMark/gastro-api/internal/auth"
	"github.com/EduardoMark/gastro-api/internal/middleware"
	"github.com/EduardoMark/gastro-api/pkg/jsonutils"
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

	body, err := jsonutils.DecodeJson[SignupRequest](r)
	if err != nil {
		jsonutils.EncodeJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if err := body.Validate(); err != nil {
		jsonutils.EncodeJson(w, http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := h.s.Create(
		ctx,
		body.Name,
		body.Email,
		body.Password,
		body.Role,
	); err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			jsonutils.EncodeJson(w, http.StatusConflict, map[string]string{
				"error": "email already exists",
			})
			return
		}

		jsonutils.EncodeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	jsonutils.EncodeJson(w, http.StatusCreated, map[string]string{
		"success": "user created with success",
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Login handler running...")
	ctx := r.Context()

	body, err := jsonutils.DecodeJson[LoginRequest](r)
	if err != nil {
		jsonutils.EncodeJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid body request",
		})
		return
	}

	if err := body.Validate(); err != nil {
		jsonutils.EncodeJson(w, http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
		return
	}

	user, err := h.s.Authenticate(ctx, body.Email, body.Password)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			jsonutils.EncodeJson(w, http.StatusNotFound, map[string]string{
				"error": "user not found",
			})
			return
		}

		if errors.Is(err, ErrInvalidCredentials) {
			jsonutils.EncodeJson(w, http.StatusBadRequest, map[string]string{
				"error": "invalid credentials",
			})
			return
		}

		jsonutils.EncodeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	token, err := h.authService.New(user.ID.String(), string(user.Role))
	if err != nil {
		jsonutils.EncodeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	jsonutils.EncodeJson(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Change Password handler running...")
	ctx := r.Context()

	userIDRaw := ctx.Value(middleware.CtxUserId)
	userID, ok := userIDRaw.(string)
	if !ok {
		jsonutils.EncodeJson(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid user id in context",
		})
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		jsonutils.EncodeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "invalid user id type uuuid",
		})
		return
	}

	body, err := jsonutils.DecodeJson[ChangePasswordRequest](r)
	if err != nil {
		jsonutils.EncodeJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid body request",
		})
		return
	}

	if err := body.Validate(); err != nil {
		jsonutils.EncodeJson(w, http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := h.s.ChangePassword(ctx, uid, body.NewPassword); err != nil {
		if errors.Is(err, ErrSamePassword) {
			jsonutils.EncodeJson(w, http.StatusBadRequest, map[string]string{
				"error": "new password cannot be the same as the old password",
			})
			return
		}

		if errors.Is(err, ErrUserNotFound) {
			jsonutils.EncodeJson(w, http.StatusNotFound, map[string]string{
				"error": "user not found",
			})
			return
		}

		jsonutils.EncodeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

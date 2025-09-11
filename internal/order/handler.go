package order

import (
	"net/http"

	"github.com/EduardoMark/gastro-api/internal/middleware"
	"github.com/EduardoMark/gastro-api/pkg/jsonutils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	s   Service
	jwt middleware.JWTMiddleware
}

func NewOrderHandler(s Service, jwt middleware.JWTMiddleware) OrderHandler {
	return OrderHandler{
		s:   s,
		jwt: jwt,
	}
}

func (h *OrderHandler) OrderRoutes(r chi.Router) {
	r.Route("/orders", func(r chi.Router) {
		r.Use(h.jwt.JWTAuth)

		r.Post("/", h.Create)
	})
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Create Order running...")

	ctx := r.Context()
	idRaw, ok := ctx.Value(middleware.CtxUserId).(string)
	if !ok {
		jsonutils.EncodeJson(w, http.StatusBadRequest, map[string]string{
			"error": "user id not found",
		})
		return
	}

	userID, err := uuid.Parse(idRaw)
	if err != nil {
		jsonutils.EncodeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "invalid user id type uuuid",
		})
		return
	}

	body, err := jsonutils.DecodeJson[CreateOrderRequest](r)
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

	if err := h.s.Create(ctx, userID, body.Items); err != nil {
		jsonutils.EncodeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	jsonutils.EncodeJson(w, http.StatusCreated, map[string]string{
		"success": "order created with success",
	})
}

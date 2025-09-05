package dishes

import (
	"errors"
	"net/http"

	"github.com/EduardoMark/gastro-api/internal/middleware"
	"github.com/EduardoMark/gastro-api/pkg/jsonutils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DishHandler struct {
	s   Service
	jwt *middleware.JWTMiddleware
}

func NewDishHandler(s Service, jwt *middleware.JWTMiddleware) DishHandler {
	return DishHandler{
		s:   s,
		jwt: jwt,
	}
}

func (h *DishHandler) DishRoutes(r chi.Router) {
	r.Route("/dishes", func(r chi.Router) {
		r.Use(h.jwt.JWTAuth)

		r.Post("/", h.Create)
		r.Get("/{id}", h.GetOne)
		r.Get("/", h.Query)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func (h *DishHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := jsonutils.DecodeJson[CreateRequest](r)
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

	if err := h.s.Create(ctx, body.Name, body.Description, body.Category, body.Price); err != nil {
		if errors.Is(err, ErrDishAlreadyExists) {
			jsonutils.EncodeJson(w, http.StatusConflict, map[string]string{
				"error": "dish already exists",
			})
			return
		}
	}

	jsonutils.EncodeJson(w, http.StatusCreated, map[string]string{
		"error": "dish created with success",
	})
}

func (h *DishHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idRaw := chi.URLParam(r, "id")

	id, err := uuid.Parse(idRaw)
	if err != nil {
		jsonutils.EncodeJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid uuid type",
		})
		return
	}

	record, err := h.s.GetOneByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrDishNotFound) {
			jsonutils.EncodeJson(w, http.StatusNotFound, map[string]string{
				"error": "dish not found",
			})
			return
		}

		jsonutils.EncodeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	response := DishResponse{
		ID:          record.ID.String(),
		Name:        record.Name,
		Description: record.Description,
		Price:       record.Price.String(),
		Category:    record.Category,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
	}

	jsonutils.EncodeJson(w, http.StatusOK, map[string]DishResponse{
		"dish": response,
	})
}

func (h *DishHandler) Query(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	records, err := h.s.Query(ctx)
	if err != nil {
		if errors.Is(err, ErrDishNotFound) {
			jsonutils.EncodeJson(w, http.StatusNotFound, map[string]string{
				"error": "dish not found",
			})
			return
		}

		jsonutils.EncodeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	response := make([]DishResponse, len(records))
	for i, record := range records {
		response[i] = DishResponse{
			ID:          record.ID.String(),
			Name:        record.Name,
			Description: record.Description,
			Price:       record.Price.String(),
			Category:    record.Category,
			CreatedAt:   record.CreatedAt,
			UpdatedAt:   record.UpdatedAt,
		}
	}

	jsonutils.EncodeJson(w, http.StatusOK, map[string][]DishResponse{
		"dishes": response,
	})
}

func (h *DishHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idRaw := chi.URLParam(r, "id")

	id, err := uuid.Parse(idRaw)
	if err != nil {
		jsonutils.EncodeJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid uuid type",
		})
		return
	}

	body, err := jsonutils.DecodeJson[UpdateRequest](r)
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

	if err := h.s.Update(ctx, id, body); err != nil {
		if errors.Is(err, ErrDishNotFound) {
			jsonutils.EncodeJson(w, http.StatusNotFound, map[string]string{
				"error": "dish not found",
			})
			return
		}

		jsonutils.EncodeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	jsonutils.EncodeJson(w, http.StatusOK, map[string]string{
		"success": "dish updated with success",
	})
}

func (h *DishHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idRaw := chi.URLParam(r, "id")

	id, err := uuid.Parse(idRaw)
	if err != nil {
		jsonutils.EncodeJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid uuid type",
		})
		return
	}

	if err := h.s.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrDishNotFound) {
			jsonutils.EncodeJson(w, http.StatusNotFound, map[string]string{
				"error": "dish not found",
			})
			return
		}

		jsonutils.EncodeJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid uuid type",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

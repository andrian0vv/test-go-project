package api_orders_post

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/andrian0vv/test-go-project/internal/api"
	ie "github.com/andrian0vv/test-go-project/internal/errors"
	"github.com/andrian0vv/test-go-project/internal/services/orders"
)

type orderService interface {
	CreateOrder(ctx context.Context, order orders.Order) error
}

type Handler struct {
	orderService orderService
	log          *slog.Logger
}

func New(orderService orderService, log *slog.Logger) *Handler {
	return &Handler{
		orderService: orderService,
		log:          log,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var req api.PostApiOrdersJSONBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		body, _ := json.Marshal(api.ToError(ie.WithCode(err, api.CodeValidation)))
		_, _ = w.Write(body)
		h.log.Error("failed to decode body", err)

		return
	}

	if err := validate(req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		body, _ := json.Marshal(api.ToError(err))
		_, _ = w.Write(body)
		h.log.Error("invalid request", err)

		return
	}

	err := h.orderService.CreateOrder(r.Context(), convertOrder(req))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		body, _ := json.Marshal(api.ToError(err))
		_, _ = w.Write(body)
		h.log.Error("failed to create order", err)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func validate(req api.PostApiOrdersJSONBody) error {
	if req.UserId <= 0 {
		return ie.WithCode(errors.New("invalid userId"), api.CodeValidation)
	}

	if len(req.Products) == 0 {
		return ie.WithCode(errors.New("empty products"), api.CodeValidation)
	}

	for _, p := range req.Products {
		if p.Quantity <= 0 {
			return ie.WithCode(errors.New("empty products"), api.CodeValidation)
		}
	}

	return nil
}

func convertOrder(req api.PostApiOrdersJSONBody) orders.Order {
	products := make([]orders.Product, 0, len(req.Products))
	for _, p := range req.Products {
		products = append(products, orders.Product{
			ProductID: p.Id,
			Quantity:  p.Quantity,
		})
	}

	return orders.Order{
		UserID:   req.UserId, // Assume that we have auth layer before service, so UserId is valid
		Products: products,
	}
}

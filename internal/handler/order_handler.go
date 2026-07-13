package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"api_clientes/internal/model"
	"api_clientes/internal/repository"
)

// OrderHandler handles HTTP requests for the /orders resource.
type OrderHandler struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
}

// NewOrderHandler creates an OrderHandler.
func NewOrderHandler(
	orderRepo *repository.OrderRepository,
	productRepo *repository.ProductRepository,
) *OrderHandler {
	return &OrderHandler{orderRepo: orderRepo, productRepo: productRepo}
}

// Create handles POST /orders.
//
// Happy path  → 201 with the created order and its items.
// Error paths → 400 (validation), 404 (customer/product not found),
//
//	422 (insufficient stock), 500.
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if req.CustomerID == "" {
		respondError(w, http.StatusBadRequest, "customer_id é obrigatório")
		return
	}
	if len(req.Items) == 0 {
		respondError(w, http.StatusBadRequest, "pelo menos um item é obrigatório")
		return
	}
	for _, item := range req.Items {
		if item.ProductID == "" {
			respondError(w, http.StatusBadRequest, "product_id é obrigatório em todos os itens")
			return
		}
		if item.Quantity <= 0 {
			respondError(w, http.StatusBadRequest, "quantity deve ser positivo em todos os itens")
			return
		}
	}

	order, err := h.orderRepo.Create(&req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			respondError(w, http.StatusNotFound, "cliente ou produto não encontrado")
		case errors.Is(err, repository.ErrInsufficientStock):
			respondError(w, http.StatusUnprocessableEntity, "estoque insuficiente para um ou mais produtos")
		default:
			respondError(w, http.StatusInternalServerError, "erro ao criar pedido")
		}
		return
	}

	respondJSON(w, http.StatusCreated, order)
}

// GetByID handles GET /orders/{id}.
func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	order, err := h.orderRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, http.StatusNotFound, "pedido não encontrado")
			return
		}
		respondError(w, http.StatusInternalServerError, "erro ao buscar pedido")
		return
	}

	respondJSON(w, http.StatusOK, order)
}

// List handles GET /orders?limit=10&offset=0.
func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := 10, 0

	if v := r.URL.Query().Get("limit"); v != "" {
		l, err := strconv.Atoi(v)
		if err != nil || l <= 0 {
			respondError(w, http.StatusBadRequest, "parâmetro limit inválido")
			return
		}
		if l > 100 {
			l = 100 // cap máximo
		}
		limit = l
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		o, err := strconv.Atoi(v)
		if err != nil || o < 0 {
			respondError(w, http.StatusBadRequest, "parâmetro offset inválido")
			return
		}
		offset = o
	}

	orders, err := h.orderRepo.List(limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "erro ao listar pedidos")
		return
	}

	if orders == nil {
		orders = []*model.Order{}
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"data":   orders,
		"limit":  limit,
		"offset": offset,
	})
}

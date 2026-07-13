package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"api_clientes/internal/model"
	"api_clientes/internal/repository"
)

// ProductHandler handles HTTP requests for the /products resource.
type ProductHandler struct {
	repo *repository.ProductRepository
}

// NewProductHandler creates a ProductHandler.
func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// Create handles POST /products.
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "name é obrigatório")
		return
	}
	if req.Price < 0 {
		respondError(w, http.StatusBadRequest, "price deve ser maior ou igual a zero")
		return
	}
	if req.StockQuantity < 0 {
		respondError(w, http.StatusBadRequest, "stock_quantity deve ser maior ou igual a zero")
		return
	}

	product, err := h.repo.Create(&req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "erro ao criar produto")
		return
	}

	respondJSON(w, http.StatusCreated, product)
}

// List handles GET /products.
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.List()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "erro ao listar produtos")
		return
	}

	if products == nil {
		products = []*model.Product{}
	}

	respondJSON(w, http.StatusOK, products)
}

// GetByID handles GET /products/{id}.
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	product, err := h.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, http.StatusNotFound, "produto não encontrado")
			return
		}
		respondError(w, http.StatusInternalServerError, "erro ao buscar produto")
		return
	}

	respondJSON(w, http.StatusOK, product)
}

// Update handles PUT /products/{id}.
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req model.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "name é obrigatório")
		return
	}
	if req.Price < 0 {
		respondError(w, http.StatusBadRequest, "price deve ser maior ou igual a zero")
		return
	}
	if req.StockQuantity < 0 {
		respondError(w, http.StatusBadRequest, "stock_quantity deve ser maior ou igual a zero")
		return
	}

	product, err := h.repo.Update(id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, http.StatusNotFound, "produto não encontrado")
			return
		}
		respondError(w, http.StatusInternalServerError, "erro ao atualizar produto")
		return
	}

	respondJSON(w, http.StatusOK, product)
}

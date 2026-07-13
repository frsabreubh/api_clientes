package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"api_clientes/internal/model"
	"api_clientes/internal/repository"
)

// CustomerHandler handles HTTP requests for the /customers resource.
type CustomerHandler struct {
	repo *repository.CustomerRepository
}

// NewCustomerHandler creates a CustomerHandler.
func NewCustomerHandler(repo *repository.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{repo: repo}
}

// Create handles POST /customers.
//
// Happy path  → 201 Created with the new customer (password_hash omitted).
// Error paths → 400 (missing fields), 409 (duplicate e-mail), 500.
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Name == "" || req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "name, email e password são obrigatórios")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "erro ao processar senha")
		return
	}

	customer, err := h.repo.Create(req.Name, req.Email, string(hash))
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			respondError(w, http.StatusConflict, "e-mail já cadastrado")
			return
		}
		respondError(w, http.StatusInternalServerError, "erro ao criar cliente")
		return
	}

	respondJSON(w, http.StatusCreated, customer)
}

// GetByID handles GET /customers/{id}.
func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	customer, err := h.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, http.StatusNotFound, "cliente não encontrado")
			return
		}
		respondError(w, http.StatusInternalServerError, "erro ao buscar cliente")
		return
	}

	respondJSON(w, http.StatusOK, customer)
}

// Login handles POST /customers/login.
//
// Happy path  → 200 with customer data.
// Error paths → 400 (missing fields), 401 (wrong credentials), 500.
func (h *CustomerHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "email e password são obrigatórios")
		return
	}

	customer, err := h.repo.GetByEmail(strings.ToLower(req.Email))
	if err != nil {
		// Return 401 for both "not found" and DB errors to avoid user enumeration.
		respondError(w, http.StatusUnauthorized, "credenciais inválidas")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(req.Password)); err != nil {
		respondError(w, http.StatusUnauthorized, "credenciais inválidas")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"message":  "login realizado com sucesso",
		"customer": customer,
	})
}

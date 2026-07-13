package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"api_clientes/internal/handler"
)

// Setup builds and returns the application router.
func Setup(
	customerHandler *handler.CustomerHandler,
	productHandler *handler.ProductHandler,
	orderHandler *handler.OrderHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Customers
	r.Route("/customers", func(r chi.Router) {
		r.Post("/", customerHandler.Create)
		r.Post("/login", customerHandler.Login)
		r.Get("/{id}", customerHandler.GetByID)
	})

	// Products
	r.Route("/products", func(r chi.Router) {
		r.Post("/", productHandler.Create)
		r.Get("/", productHandler.List)
		r.Get("/{id}", productHandler.GetByID)
		r.Put("/{id}", productHandler.Update)
	})

	// Orders
	r.Route("/orders", func(r chi.Router) {
		r.Post("/", orderHandler.Create)
		r.Get("/", orderHandler.List)
		r.Get("/{id}", orderHandler.GetByID)
	})

	return r
}

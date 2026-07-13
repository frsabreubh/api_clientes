package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"api_clientes/internal/config"
	"api_clientes/internal/db"
	"api_clientes/internal/handler"
	"api_clientes/internal/migrate"
	"api_clientes/internal/repository"
	"api_clientes/internal/router"
)

func main() {
	// Load .env (optional; falls back to real env vars).
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	cfg := config.Load()

	// Connect to PostgreSQL.
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}
	defer database.Close()

	// Run pending SQL migrations.
	if err := migrate.Run(database); err != nil {
		log.Fatalf("Falha ao executar migrations: %v", err)
	}

	// Wire up repositories.
	customerRepo := repository.NewCustomerRepository(database)
	productRepo := repository.NewProductRepository(database)
	orderRepo := repository.NewOrderRepository(database)

	// Wire up handlers.
	customerHandler := handler.NewCustomerHandler(customerRepo)
	productHandler := handler.NewProductHandler(productRepo)
	orderHandler := handler.NewOrderHandler(orderRepo, productRepo)

	// Build router.
	r := router.Setup(customerHandler, productHandler, orderHandler)

	addr := ":" + cfg.Port
	fmt.Printf("API rodando em http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

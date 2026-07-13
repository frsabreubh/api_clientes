package repository

import (
	"database/sql"
	"errors"

	"api_clientes/internal/model"
)

// ProductRepository handles all product SQL operations.
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a ProductRepository backed by db.
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create inserts a new product.
func (r *ProductRepository) Create(req *model.CreateProductRequest) (*model.Product, error) {
	const q = `
		INSERT INTO products (name, description, price, stock_quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, price, stock_quantity, created_at, updated_at`

	p := &model.Product{}
	err := r.db.QueryRow(q, req.Name, req.Description, req.Price, req.StockQuantity).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetByID fetches a product by primary key.
func (r *ProductRepository) GetByID(id string) (*model.Product, error) {
	const q = `
		SELECT id, name, description, price, stock_quantity, created_at, updated_at
		FROM products
		WHERE id = $1`

	p := &model.Product{}
	err := r.db.QueryRow(q, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return p, nil
}

// List returns all products ordered by creation date.
func (r *ProductRepository) List() ([]*model.Product, error) {
	const q = `
		SELECT id, name, description, price, stock_quantity, created_at, updated_at
		FROM products
		ORDER BY created_at DESC`

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		p := &model.Product{}
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Price,
			&p.StockQuantity, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// Update replaces all mutable fields of a product.
func (r *ProductRepository) Update(id string, req *model.UpdateProductRequest) (*model.Product, error) {
	const q = `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock_quantity = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, name, description, price, stock_quantity, created_at, updated_at`

	p := &model.Product{}
	err := r.db.QueryRow(q, req.Name, req.Description, req.Price, req.StockQuantity, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return p, nil
}

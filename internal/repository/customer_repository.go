package repository

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"

	"api_clientes/internal/model"
)

// CustomerRepository handles all customer SQL operations.
type CustomerRepository struct {
	db *sql.DB
}

// NewCustomerRepository creates a CustomerRepository backed by db.
func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// Create inserts a new customer and returns the persisted record.
func (r *CustomerRepository) Create(name, email, passwordHash string) (*model.Customer, error) {
	const q = `
		INSERT INTO customers (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, password_hash, created_at, updated_at`

	c := &model.Customer{}
	err := r.db.QueryRow(q, name, email, passwordHash).Scan(
		&c.ID, &c.Name, &c.Email, &c.PasswordHash, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, ErrDuplicate
		}
		return nil, err
	}
	return c, nil
}

// GetByID fetches a customer by primary key.
func (r *CustomerRepository) GetByID(id string) (*model.Customer, error) {
	const q = `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM customers
		WHERE id = $1`

	c := &model.Customer{}
	err := r.db.QueryRow(q, id).Scan(
		&c.ID, &c.Name, &c.Email, &c.PasswordHash, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

// GetByEmail fetches a customer by e-mail (used for login).
func (r *CustomerRepository) GetByEmail(email string) (*model.Customer, error) {
	const q = `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM customers
		WHERE email = $1`

	c := &model.Customer{}
	err := r.db.QueryRow(q, email).Scan(
		&c.ID, &c.Name, &c.Email, &c.PasswordHash, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

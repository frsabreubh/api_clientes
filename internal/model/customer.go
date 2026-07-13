package model

import "time"

// Customer represents a registered user.
type Customer struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // never serialised
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateCustomerRequest is the payload for POST /customers.
type CreateCustomerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest is the payload for POST /customers/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

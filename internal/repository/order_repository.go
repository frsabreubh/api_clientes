package repository

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"

	"api_clientes/internal/model"
)

// OrderRepository handles all order SQL operations.
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates an OrderRepository backed by db.
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create opens a transaction and:
//  1. Locks each requested product row (SELECT … FOR UPDATE).
//  2. Validates stock availability.
//  3. Inserts the order.
//  4. Inserts every order_item and decrements stock atomically.
//
// The entire operation is rolled back on any error.
func (r *OrderRepository) Create(req *model.CreateOrderRequest) (*model.Order, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	order, err := r.createTx(tx, req)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return order, nil
}

// createTx runs the transactional order-creation logic.
func (r *OrderRepository) createTx(tx *sql.Tx, req *model.CreateOrderRequest) (*model.Order, error) {
	type productInfo struct {
		price float64
		stock int
	}

	// 1. Lock products and validate stock.
	infos := make(map[string]productInfo, len(req.Items))
	for _, item := range req.Items {
		var info productInfo
		err := tx.QueryRow(`
			SELECT price, stock_quantity
			FROM products
			WHERE id = $1
			FOR UPDATE`, item.ProductID,
		).Scan(&info.price, &info.stock)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		if info.stock < item.Quantity {
			return nil, ErrInsufficientStock
		}

		infos[item.ProductID] = info
	}

	// 2. Calculate total.
	var total float64
	for _, item := range req.Items {
		total += infos[item.ProductID].price * float64(item.Quantity)
	}

	// 3. Insert order.
	order := &model.Order{}
	err := tx.QueryRow(`
		INSERT INTO orders (customer_id, status, total_amount)
		VALUES ($1, 'pending', $2)
		RETURNING id, customer_id, status, total_amount, created_at, updated_at`,
		req.CustomerID, total,
	).Scan(
		&order.ID, &order.CustomerID, &order.Status,
		&order.TotalAmount, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			// foreign_key_violation → customer not found
			return nil, ErrNotFound
		}
		return nil, err
	}

	// 4. Insert items and decrement stock.
	for _, item := range req.Items {
		unitPrice := infos[item.ProductID].price

		// Decrement stock (check again inside the transaction).
		res, err := tx.Exec(`
			UPDATE products
			SET stock_quantity = stock_quantity - $1, updated_at = NOW()
			WHERE id = $2 AND stock_quantity >= $1`,
			item.Quantity, item.ProductID,
		)
		if err != nil {
			return nil, err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		if affected == 0 {
			return nil, ErrInsufficientStock
		}

		// Insert order_item.
		var oi model.OrderItem
		err = tx.QueryRow(`
			INSERT INTO order_items (order_id, product_id, quantity, unit_price)
			VALUES ($1, $2, $3, $4)
			RETURNING id, order_id, product_id, quantity, unit_price, created_at`,
			order.ID, item.ProductID, item.Quantity, unitPrice,
		).Scan(
			&oi.ID, &oi.OrderID, &oi.ProductID,
			&oi.Quantity, &oi.UnitPrice, &oi.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		order.Items = append(order.Items, oi)
	}

	return order, nil
}

// GetByID fetches an order and its items.
func (r *OrderRepository) GetByID(id string) (*model.Order, error) {
	order := &model.Order{}
	err := r.db.QueryRow(`
		SELECT id, customer_id, status, total_amount, created_at, updated_at
		FROM orders
		WHERE id = $1`, id,
	).Scan(
		&order.ID, &order.CustomerID, &order.Status,
		&order.TotalAmount, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	rows, err := r.db.Query(`
		SELECT id, order_id, product_id, quantity, unit_price, created_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var oi model.OrderItem
		if err := rows.Scan(
			&oi.ID, &oi.OrderID, &oi.ProductID,
			&oi.Quantity, &oi.UnitPrice, &oi.CreatedAt,
		); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, oi)
	}
	return order, rows.Err()
}

// List returns orders with pagination (limit / offset), newest first.
func (r *OrderRepository) List(limit, offset int) ([]*model.Order, error) {
	rows, err := r.db.Query(`
		SELECT id, customer_id, status, total_amount, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		o := &model.Order{}
		if err := rows.Scan(
			&o.ID, &o.CustomerID, &o.Status,
			&o.TotalAmount, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

package repository

import (
	"database/sql"
	"fmt"
	"time"

	"lesson23/model"

	"github.com/google/uuid"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order *model.Order) error {
	id := uuid.NewString()

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Commit()

	query := `
		INSERT INTO orders (id, customer_id, products, total_amount, status, shipping_address, payment_method, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = tx.Exec(query, id, order.CustomerID, order.Products, order.TotalAmount, order.Status, order.ShippingAddress, order.PaymentMethod, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error inserting order: %w", err)
	}

	return nil
}

func (r *OrderRepository) GetOrderById(id string) (*model.Order, error) {
	var order model.Order
	query := `
		SELECT id, customer_id, products, total_amount, status, shipping_address, payment_method, created_at, updated_at
		FROM orders WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.CustomerID,
		&order.Products,
		&order.TotalAmount,
		&order.Status,
		&order.ShippingAddress,
		&order.PaymentMethod,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order with id '%s' not found", id)
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving order: %w", err)
	}

	return &order, nil
}

func (r *OrderRepository) GetAllOrders() ([]model.Order, error) {
	query := `
		SELECT id, customer_id, products, total_amount, status, shipping_address, payment_method, created_at, updated_at
		FROM orders`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error retrieving orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.Products,
			&order.TotalAmount,
			&order.Status,
			&order.ShippingAddress,
			&order.PaymentMethod,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) UpdateOrder(id string, order *model.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Commit()

	var existingOrder model.Order
	err = tx.QueryRow("SELECT id FROM orders WHERE id = $1", id).Scan(&existingOrder.ID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("order with id '%s' not found", id)
	} else if err != nil {
		return fmt.Errorf("error checking if order exists: %w", err)
	}

	query := `
		UPDATE orders 
		SET customer_id = $1, products = $2, total_amount = $3, status = $4, 
			shipping_address = $5, payment_method = $6, updated_at = $7
		WHERE id = $8`
	_, err = tx.Exec(query, order.CustomerID, order.Products, order.TotalAmount, order.Status, order.ShippingAddress, order.PaymentMethod, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error updating order: %w", err)
	}

	return nil
}

func (r *OrderRepository) DeleteOrder(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Commit()

	var existingOrder model.Order
	err = tx.QueryRow("SELECT id FROM orders WHERE id = $1", id).Scan(&existingOrder.ID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("order with id '%s' not found", id)
	} else if err != nil {
		return fmt.Errorf("error checking if order exists: %w", err)
	}

	_, err = tx.Exec("DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting order: %w", err)
	}

	return nil
}

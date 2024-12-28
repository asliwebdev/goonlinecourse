package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"lesson22.5/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(product *model.Product) error {
	id := uuid.NewString()

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	_, err = tx.Exec("INSERT INTO product (id, name, description, price, category, stock, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		id, product.Name, product.Description, product.Price, product.Category, product.Stock, product.CreatedAt, product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error creating product: %w", err)
	}

	return nil
}

func (r *ProductRepository) GetAllProducts() ([]model.Product, error) {
	rows, err := r.db.Query("SELECT id, name, description, price, category, stock, created_at, updated_at FROM product")
	if err != nil {
		return nil, fmt.Errorf("error fetching products: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.Stock, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning product: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through products: %w", err)
	}

	return products, nil
}

func (r *ProductRepository) GetProductById(id string) (*model.Product, error) {
	var product model.Product
	err := r.db.QueryRow("SELECT id, name, description, price, category, stock, created_at, updated_at FROM product WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found: %w", err)
		}
		return nil, fmt.Errorf("error fetching product by ID: %w", err)
	}

	return &product, nil
}

func (r *ProductRepository) UpdateProduct(id string, product *model.Product) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Commit()

	var existingProduct model.Product
	err = tx.QueryRow("SELECT id FROM product WHERE id = $1", id).Scan(&existingProduct.ID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("product with id '%s' not found", id)
	} else if err != nil {
		return fmt.Errorf("error checking if product exists: %w", err)
	}

	_, err = tx.Exec("UPDATE product SET name = $1, description = $2, price = $3, category = $4, stock = $5, updated_at = $6 WHERE id = $7",
		product.Name, product.Description, product.Price, product.Category, product.Stock, product.UpdatedAt, id)
	if err != nil {
		return fmt.Errorf("error updating product: %w", err)
	}

	return nil
}

func (r *ProductRepository) DeleteProduct(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Commit()

	var existingProduct model.Product
	err = tx.QueryRow("SELECT id FROM product WHERE id = $1", id).Scan(&existingProduct.ID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("product with id '%s' not found", id)
	} else if err != nil {
		return fmt.Errorf("error checking if product exists: %w", err)
	}

	_, err = tx.Exec("DELETE FROM product WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}

	return nil
}

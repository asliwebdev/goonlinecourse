package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"lesson22.5/model"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(category *model.Category) error {
	var existingCategory model.Category
	err := r.db.QueryRow("SELECT id, name FROM category WHERE name = $1", category.Name).Scan(&existingCategory.ID, &existingCategory.Name)

	if err == nil {
		return fmt.Errorf("category with name '%s' already exists", category.Name)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("error checking if category exists: %w", err)
	}

	id := uuid.NewString()

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	_, err = tx.Exec("insert into category (id, name) values ($1, $2)", id, category.Name)
	if err != nil {
		return fmt.Errorf("error creating category: %w", err)
	}

	return nil
}

func (r *CategoryRepository) GetCategoryById(id string) (*model.Category, error) {
	category := &model.Category{}
	err := r.db.QueryRow("SELECT * FROM category WHERE id = $1", id).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching category by ID: %w", err)
	}

	return category, nil
}

func (r *CategoryRepository) GetAllCategories() ([]model.Category, error) {
	rows, err := r.db.Query("SELECT id, name, created_at, updated_at FROM category")
	if err != nil {
		return nil, fmt.Errorf("error fetching categories: %w", err)
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning category row: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows: %w", err)
	}

	return categories, nil
}

func (r *CategoryRepository) UpdateCategory(id string, category *model.Category) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Commit()

	var existingCategory model.Category
	err = tx.QueryRow("SELECT id FROM category WHERE id = $1", id).Scan(&existingCategory.ID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("category with id '%s' not found", id)
	} else if err != nil {
		return fmt.Errorf("error checking if category exists: %w", err)
	}

	_, err = tx.Exec("UPDATE category SET name = $1 WHERE id = $2", category.Name, id)
	if err != nil {
		return fmt.Errorf("error updating category: %w", err)
	}

	return nil
}

func (r *CategoryRepository) DeleteCategory(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Commit()

	var existingCategory model.Category
	err = tx.QueryRow("SELECT id FROM category WHERE id = $1", id).Scan(&existingCategory.ID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("category with id '%s' not found", id)
	} else if err != nil {
		return fmt.Errorf("error checking if category exists: %w", err)
	}

	_, err = tx.Exec("DELETE FROM category WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting category: %w", err)
	}

	return nil
}

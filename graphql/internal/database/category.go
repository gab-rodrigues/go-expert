package database

import (
	"database/sql"
	"github.com/google/uuid"
)

type Category struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
}

func NewCategory(db *sql.DB) *Category {
	return &Category{
		db: db,
	}
}

func (c *Category) Create(name, description string) (Category, error) {
	id := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO categories (id, name, description) VALUES (?, ?, ?)", id, name, description)
	if err != nil {
		return Category{}, err
	}
	return Category{
		ID:          id,
		Name:        name,
		Description: description,
	}, nil
}

func (c *Category) FindAll() ([]Category, error) {
	rows, err := c.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *Category) FindByCourseID(courseID string) (Category, error) {
	row := c.db.QueryRow(`
		SELECT categories.id, categories.name, categories.description
		FROM categories
		JOIN courses ON categories.id = courses.category_id
		WHERE courses.id = ?`, courseID)

	var category Category
	if err := row.Scan(&category.ID, &category.Name, &category.Description); err != nil {
		if err == sql.ErrNoRows {
			return Category{}, nil // No category found for the given course ID
		}
		return Category{}, err
	}

	return category, nil
}

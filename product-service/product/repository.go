package product

import (
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(input ProductCreate) (Product, error) {
	result, err := r.db.Exec(
		"INSERT INTO products (name) VALUES (?)",
		input.Name,
	)
	if err != nil {
		return Product{}, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return Product{}, err
	}

	return Product{
		ID:   int(id),
		Name: input.Name,
	}, nil
}

func (r *Repository) GetByID(id int) (Product, error) {
	var product Product

	err := r.db.QueryRow(
		"SELECT id, name FROM products WHERE id = ?",
		id,
	).Scan(
		&product.ID,
		&product.Name,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Product{}, ErrProductNotFound
		}

		return Product{}, err
	}

	return product, nil
}

func (r *Repository) List() ([]Product, error) {

	rows, err := r.db.Query(
		"SELECT id, name FROM products",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product

		if err := rows.Scan(
			&product.ID,
			&product.Name,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil

}

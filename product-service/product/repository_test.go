package product

import (
	"database/sql"
	"errors"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestRepository_CreateAndGetByID(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	input := ProductCreate{
		Name: "Coffee",
	}

	// Act
	createdProduct, err := repo.Create(input)
	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetByID(createdProduct.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if got.ID != createdProduct.ID {
		t.Errorf("expected ID %d, got %d", createdProduct.ID, got.ID)
	}

	if got.Name != "Coffee" {
		t.Errorf("expected name Coffee, got %s", got.Name)
	}
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	// Act
	_, err := repo.GetByID(999)

	// Assert
	if !errors.Is(err, ErrProductNotFound) {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestRepository_List(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	_, err := repo.Create(ProductCreate{Name: "Coffee"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Create(ProductCreate{Name: "Tea"})
	if err != nil {
		t.Fatal(err)
	}

	// Act
	products, err := repo.List()
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}

	if products[0].Name != "Coffee" {
		t.Errorf("expected first product Coffee, got %s", products[0].Name)
	}

	if products[1].Name != "Tea" {
		t.Errorf("expected second product Tea, got %s", products[1].Name)
	}
}

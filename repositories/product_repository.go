package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"kasir-go/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) FindAll(name string) ([]models.Product, error) {
	query := "SELECT id, name, price, stock, category_id FROM products"

	var args []interface{}
	if name != "" {
		query += " WHERE name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	query += " ORDER BY created_at DESC"

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"

	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)

	return err
}

func (repo *ProductRepository) FindById(id int) (*models.Product, error) {
	query := "SELECT id, name, price, stock, category_id FROM products WHERE id = $1"

	var product models.Product
	err := repo.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("product id %d not found", id)
	}

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"

	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"

	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return err
}

func (repo *ProductRepository) FindByCategoryId(categoryId int) ([]models.Product, error) {
	query := "SELECT id, name, price, stock FROM products where category_id = $1"

	rows, err := repo.db.Query(query, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

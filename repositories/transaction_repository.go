package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"kasir-go/models"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	var (
		res *models.Transaction
	)

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0

	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productName string
		var productID, price, stock int

		err := tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productID, &productName, &price, &stock)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}

		subtotal := item.Quantity * price
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, productID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		_, err := tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)", details[i].TransactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	res = &models.Transaction{
		ID:                 transactionID,
		TotalAmount:        totalAmount,
		TransactionDetails: details,
	}

	return res, nil
}

func (r *TransactionRepository) GetSummaryByPeriod(start, end time.Time) (totalRevenue int, totalTransaction int, err error) {
	query := `
		SELECT
			COALESCE(SUM(total_amount), 0) AS total_revenue,
			COUNT(*) AS total_transaction
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
	`

	err = r.db.QueryRow(query, start, end).Scan(&totalRevenue, &totalTransaction)
	if err != nil {
		return 0, 0, err
	}

	return totalRevenue, totalTransaction, nil
}

func (r *TransactionRepository) GetBestSellingProductByPeriod(start, end time.Time) (name string, quantity int, err error) {
	query := `
		SELECT
			p.name,
			SUM(td.quantity) AS qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.name
		ORDER BY qty DESC
		LIMIT 1
	`

	err = r.db.QueryRow(query, start, end).Scan(&name, &quantity)
	if err == sql.ErrNoRows {
		return "", 0, nil
	}

	if err != nil {
		return "", 0, err
	}

	return name, quantity, nil
}

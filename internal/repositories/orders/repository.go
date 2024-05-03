package orders

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	db "github.com/andrian0vv/test-go-project/internal/database"
	"github.com/andrian0vv/test-go-project/internal/models"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateOrder(ctx context.Context, tx *sqlx.Tx, order models.Order) error {
	query := `insert into orders (user_id) values ($1) returning id`

	var orderID int64
	row := db.ChooseQuerier(tx, r.db).QueryRowxContext(ctx, query, order.UserID)
	if row.Err() != nil {
		return fmt.Errorf("query: %w", row.Err())
	}

	if err := row.Scan(&orderID); err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	for i := range order.Products {
		order.Products[i].OrderID = orderID
	}

	query = `insert into order_products (price, order_id, product_id, quantity) 
values (:price, :order_id, :product_id, :quantity)`

	_, err := db.ChooseQuerier(tx, r.db).NamedExecContext(ctx, query, order.Products)
	if err != nil {
		return fmt.Errorf("insert order products: %w", err)
	}

	return nil
}

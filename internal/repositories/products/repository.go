package products

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

func (r *Repository) GetProducts(ctx context.Context, ids []int64) (map[int64]models.Product, error) {
	query := `select id, description, price, quantity, created_at from products where id in (?)`

	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, err
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	products := make([]models.Product, 0, len(ids))

	err = r.db.SelectContext(ctx, &products, query, args...)
	if err != nil {
		return nil, err
	}

	result := make(map[int64]models.Product, len(products))
	for _, p := range products {
		result[p.ID] = p
	}

	return result, nil
}

func (r *Repository) WriteOff(ctx context.Context, tx *sqlx.Tx, products map[int64]int) error {
	query := "update products set quantity = quantity - $1 where id = $2"

	for productID, quantity := range products {
		_, err := db.ChooseQuerier(tx, r.db).ExecContext(ctx, query, quantity, productID)
		if err != nil {
			return fmt.Errorf("update product quantity: %w", err)
		}
	}

	return nil
}

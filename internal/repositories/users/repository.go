package users

import (
	"context"

	"github.com/jmoiron/sqlx"

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

func (r *Repository) CreateUser(ctx context.Context, user models.User) error {
	query := `insert into users (first_name, last_name, full_name, age, is_married, password) 
values (:first_name, :last_name, :full_name, :age, :is_married, :password)`

	_, err := r.db.NamedExecContext(ctx, query, user)

	return err
}

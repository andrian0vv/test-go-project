package models

import "time"

type User struct {
	ID        int64     `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	FullName  string    `db:"full_name"`
	Age       int       `db:"age"`
	IsMarried bool      `db:"is_married"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

type Product struct {
	ID          int64     `db:"id"`
	Description string    `db:"description"`
	Price       int64     `db:"price"`
	Tags        []string  `db:"tags"`
	Quantity    int       `db:"quantity"`
	CreatedAt   time.Time `db:"created_at"`
}

type Order struct {
	ID        int64 `db:"id"`
	UserID    int64 `db:"user_id"`
	Products  []OrderProduct
	CreatedAt time.Time `db:"created_at"`
}

type OrderProduct struct {
	ID        int64 `db:"id"`
	OrderID   int64 `db:"order_id"`
	ProductID int64 `db:"product_id"`
	Price     int64 `db:"price"`
	Quantity  int   `db:"quantity"`
}

package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"

	apiOrdersPost "github.com/andrian0vv/test-go-project/internal/api/handlers/api_orders_post"
	apiUsersPost "github.com/andrian0vv/test-go-project/internal/api/handlers/api_users_post"
	"github.com/andrian0vv/test-go-project/internal/api/middleware"
	"github.com/andrian0vv/test-go-project/internal/database"
	ordersRepository "github.com/andrian0vv/test-go-project/internal/repositories/orders"
	productsRepository "github.com/andrian0vv/test-go-project/internal/repositories/products"
	usersRepository "github.com/andrian0vv/test-go-project/internal/repositories/users"
	"github.com/andrian0vv/test-go-project/internal/services/hasher"
	"github.com/andrian0vv/test-go-project/internal/services/orders"
	"github.com/andrian0vv/test-go-project/internal/services/users"
)

func main() {
	mux := http.NewServeMux()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	db, closeDB, err := dbConnect()
	if err != nil {
		log.Error("failed to connect to database", err)
		return
	}
	defer closeDB()

	storage := database.New(db)
	ordersRepo := ordersRepository.New(db)
	productsRepo := productsRepository.New(db)
	usersRepo := usersRepository.New(db)

	userService := users.New(usersRepo, hasher.New())
	orderService := orders.New(storage, ordersRepo, productsRepo)

	apiUsersPostHandler := apiUsersPost.New(userService, log)
	apiOrdersPostHandler := apiOrdersPost.New(orderService, log)

	mux.HandleFunc("POST /api/users", apiUsersPostHandler.Handle)
	mux.HandleFunc("POST /api/orders", apiOrdersPostHandler.Handle)

	middlewares := []func(http.Handler) http.Handler{
		middleware.Metrics,
		middleware.Logging,
		middleware.Panics,
	}

	server := applyMiddlewares(mux, middlewares...)

	log.Info("Server starting on port 80...")

	err = http.ListenAndServe(":80", server)
	if err != nil {
		log.Error("Listen and serve error", err)
	}
}

func applyMiddlewares(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}

	return h
}

func dbConnect() (*sqlx.DB, func(), error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("PGUSER", "postgres"),
		getEnv("PGPASSWORD", ""),
		getEnv("PGHOST", "db"),
		getEnv("PGPORT", "5432"),
		getEnv("PGDATABASE", "master"),
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping: %w", err)
	}

	return db, func() { _ = db.Close() }, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

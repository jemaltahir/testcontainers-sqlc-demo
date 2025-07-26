package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jemaltahir/testcontainers-sqlc-demo/internal/customer"
)

func main() {
	// In real prod code read DATABASE_URL from env.
	dsn := "postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable"
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	repo := customer.NewRepository(pool)

	c, err := repo.Create(ctx, "Bob", "bob@example.com")
	if err != nil {
		log.Fatalf("insert: %v", err)
	}
	fmt.Printf("Inserted: %+v\n", c)
}

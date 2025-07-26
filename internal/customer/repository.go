package customer

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jemaltahir/testcontainers-sqlc-demo/internal/db"
)

type Repository struct{ q *db.Queries }

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{q: db.New(pool)}
}

func (r *Repository) Create(ctx context.Context, name, email string) (Customer, error) {
	row, err := r.q.CreateCustomer(ctx, db.CreateCustomerParams{Name: name, Email: email})
	if err != nil {
		return Customer{}, err
	}
	return toDomain(row), nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (Customer, error) {
	row, err := r.q.GetCustomerByEmail(ctx, email)
	if err != nil {
		return Customer{}, err
	}
	return toDomain(row), nil
}

func toDomain(c db.Customer) Customer {
	return Customer{ID: c.ID, Name: c.Name, Email: c.Email}
}

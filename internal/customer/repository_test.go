package customer

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jemaltahir/testcontainers-sqlc-demo/testhelpers"
)

var (
	pool *pgxpool.Pool
	pg   *testhelpers.Postgres
)

// TestMain runs once per package.
func TestMain(m *testing.M) {
	ctx := context.Background()

	// Reuse container locally for speed if TC_REUSE=1
	var err error
	pg, err = testhelpers.StartPostgres(ctx, os.Getenv("TC_REUSE") == "1")
	if err != nil {
		panic(err)
	}

	pool, err = pgxpool.New(ctx, pg.ConnString)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	pool.Close()
	if os.Getenv("TC_REUSE") != "1" {
		_ = pg.Terminate(ctx)
	}
	os.Exit(code)
}

func TestCreateAndGet(t *testing.T) {
	t.Parallel()
	repo := NewRepository(pool)
	ctx := context.Background()

	want, err := repo.Create(ctx, "Alice", "alice@example.com")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := repo.GetByEmail(ctx, want.Email)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != want {
		t.Fatalf("mismatch: %+v != %+v", got, want)
	}
}

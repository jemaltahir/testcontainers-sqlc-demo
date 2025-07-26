// testhelpers/postgres.go
package testhelpers

import (
	"context"
	_ "embed" // for go:embed
	"fmt"
	"os"
	"path/filepath"

	migration "github.com/jemaltahir/testcontainers-sqlc-demo/db/migration"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// -----------------------------------------------------------------------------
// Postgres helper
// -----------------------------------------------------------------------------

type Postgres struct {
	*postgres.PostgresContainer
	ConnString string
	tmpDir     string // so we can clean up the temp file
}

// StartPostgres boots a Postgres 16 container.
// Set reuse=true (TC_REUSE=1) to keep it between local runs.
func StartPostgres(ctx context.Context, reuse bool) (*Postgres, error) {
	// 1. write the embedded SQL to a temp dir
	tmpDir, err := os.MkdirTemp("", "tc-migration-*")
	if err != nil {
		return nil, fmt.Errorf("temp dir: %w", err)
	}
	sqlPath := filepath.Join(tmpDir, "000_init.sql")
	if err := os.WriteFile(sqlPath, []byte(migration.InitSQL), 0o644); err != nil {
		_ = os.RemoveAll(tmpDir)
		return nil, fmt.Errorf("write sql: %w", err)
	}

	// 2. configure container
	opts := []testcontainers.ContainerCustomizer{
		postgres.WithDatabase("appdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithInitScripts(sqlPath), // <- temp file that always exists
		postgres.BasicWaitStrategies(),
	}
	if reuse {
		opts = append(opts, testcontainers.WithReuseByName("tc-appdb"))
	}

	ctr, err := postgres.Run(ctx, "postgres:16-alpine", opts...)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return nil, fmt.Errorf("run container: %w", err)
	}

	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = testcontainers.TerminateContainer(ctr)
		_ = os.RemoveAll(tmpDir)
		return nil, fmt.Errorf("conn string: %w", err)
	}

	return &Postgres{
		PostgresContainer: ctr,
		ConnString:        dsn,
		tmpDir:            tmpDir,
	}, nil
}

// Terminate stops the container and cleans up the temp dir.
func (p *Postgres) Terminate(ctx context.Context) error {
	_ = os.RemoveAll(p.tmpDir)
	return testcontainers.TerminateContainer(p.PostgresContainer)
}

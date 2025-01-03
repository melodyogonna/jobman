package migration_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/melodyogonna/jobman/internal/migration"
)

var dbURL = "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"

func TestMigrationIsInitialized(t *testing.T) {
	err := migration.Init(dbURL)
	if err != nil {
		t.Errorf("Unable to init migration %s", err)
	}

	tableVerificationQuery := `SELECT EXISTS(SELECT FROM pg_tables WHERE schemaname = 'jobman_L1nsxfVZfj' AND tablename='jobman') as e`
	var tableExist bool
	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		t.Fatal(err)
	}

	err = db.QueryRow(context.Background(), tableVerificationQuery).Scan(&tableExist)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tableExist)
	if !tableExist {
		t.Error("Table does not exist")
	}
}

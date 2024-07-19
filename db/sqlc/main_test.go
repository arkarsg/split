package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:5432/split_app?sslmode=disable"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	testDbContainer := SetUpTestDatabase()
	defer testDbContainer.TearDown()
	testDb = testDbContainer.DbInstance

	testQueries = New(testDb)

	os.Exit(m.Run())
}

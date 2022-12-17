package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	driver   = "postgres"
	dbSource = "postgresql://root:1234@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDb, err = sql.Open(driver, dbSource)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	testQueries = New(testDb)
	os.Exit(m.Run())
}

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
var conn *sql.DB

func TestMain(m *testing.M) {
	var err error
	conn, err = sql.Open(driver, dbSource)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}

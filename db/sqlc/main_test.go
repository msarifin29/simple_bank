package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgresql"
	dbSource = "postgresql://root:12345@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	connPool, err := pgxpool.New(context.Background(), dbSource)
	//conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db : ", err)
	}
	testQueries = New(connPool)
	db := testQueries.db
	if db == nil {
		log.Fatalf("error..")
	}
	fmt.Println("good! test pass")
	os.Exit(m.Run())
}

package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/msarifin29/simple_bank/api"
	db "github.com/msarifin29/simple_bank/db/sqlc"
)

const (
	dbDriver      = "postgresql"
	dbSource      = "postgresql://root:12345@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	connPool, err := pgxpool.New(context.Background(), dbSource)
	//conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db : ", err)
	}
	store := db.NewStore(connPool)
	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}

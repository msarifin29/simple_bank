package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}
type SQLStore struct {
	connPoll *pgxpool.Pool
	*Queries
}

func NewStore(connPoll *pgxpool.Pool) Store {
	return &SQLStore{connPoll: connPoll, Queries: New(connPoll)}
}

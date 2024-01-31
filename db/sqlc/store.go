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
	*Queries
	connPoll *pgxpool.Pool
}

func NewStore(connPoll *pgxpool.Pool) *SQLStore {
	return &SQLStore{connPoll: connPoll, Queries: New(connPoll)}
}

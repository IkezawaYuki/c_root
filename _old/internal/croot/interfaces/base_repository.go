package interfaces

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/croot/infrastructre"
)

type BaseRepository interface {
	BeginTransaction(ctx context.Context) (infrastructre.Tx, error)
}

type baseRepository struct {
	dbClient infrastructre.DbClient
}

func NewBaseRepository(client infrastructre.DbClient) BaseRepository {
	return &baseRepository{
		dbClient: client,
	}
}

func (b *baseRepository) BeginTransaction(ctx context.Context) (infrastructre.Tx, error) {
	return b.dbClient.BeginTx(ctx)
}

package interfaces

import (
	"context"
	"github.com/IkezawaYuki/c_root/internal/croot/infrastructre"
)

type repository struct {
	client infrastructre.DbClient
}

type Repository interface {
	Begin(ctx context.Context) (infrastructre.Tx, error)
}

func NewRepository(client infrastructre.DbClient) Repository {
	return &repository{client: client}
}

func (r *repository) Begin(ctx context.Context) (infrastructre.Tx, error) {
	return r.client.BeginTx(ctx)
}

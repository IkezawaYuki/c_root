package infrastructre

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type DbClient interface {
	Exec(query string, args ...interface{}) (Result, error)
	Query(query string, args ...interface{}) (Rows, error)
	QueryRow(query string, args ...interface{}) Row
	Begin() (Tx, error)
	Close() error
	BeginTx(ctx context.Context) (Tx, error)
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Rows interface {
	Scan(...interface{}) error
	Next() bool
	Close() error
}

type Row interface {
	Scan(...interface{}) error
}

type Stmt interface {
	Exec(query string, args ...interface{}) (Result, error)
	Query(query string, args ...interface{}) (Row, error)
	QueryRow(query string, args ...interface{}) Row
	Close() error
}

type Tx interface {
	Commit() error
	Rollback() error
	Exec(query string, args ...interface{})
	Query(query string, args ...interface{}) (Row, error)
	QueryRow(query string, args ...interface{}) Row
}

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, ttl int) error
	Del(ctx context.Context, key string) error
	GetClient() *redis.Client
}

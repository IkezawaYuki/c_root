package infrastructre

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/IkezawaYuki/popple/config"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type mysqlClient struct {
	db *sql.DB
}

func NewMySQLClient(db *sql.DB) DbClient {
	return &mysqlClient{
		db: db,
	}
}

func GetMySQLConnection() *sql.DB {
	conn, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true",
			config.Env.DatabaseHost,
			config.Env.DatabasePass,
			config.Env.DatabaseHost,
			config.Env.DatabaseName,
		))
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}

func (m *mysqlClient) Exec(query string, args ...interface{}) (Result, error) {
	return m.db.Exec(query, args...)
}

func (m *mysqlClient) Query(query string, args ...interface{}) (Rows, error) {
	return m.db.Query(query, args...)
}

func (m *mysqlClient) QueryRow(query string, args ...interface{}) Row {
	return m.db.QueryRow(query, args...)
}

func (m *mysqlClient) QueryRowContext(ctx context.Context, query string, args ...any) Row {
	return m.db.QueryRowContext(ctx, query, args...)
}

func (m *mysqlClient) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	return m.db.QueryContext(ctx, query, args...)
}

func (m *mysqlClient) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

func (m *mysqlClient) BeginTx(ctx context.Context) (Tx, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	return Tx{
		tx: tx,
	}, err
}

func (m *mysqlClient) Close() error {
	return m.db.Close()
}

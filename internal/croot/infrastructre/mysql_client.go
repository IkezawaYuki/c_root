package infrastructre

import (
	"database/sql"
	"fmt"
	"github.com/IkezawaYuki/c_root/config"
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

func (m *mysqlClient) Begin() (Tx, error) {
	tx, err := m.db.Begin()
	t := Tx{
		tx: tx,
	}
	return t, err
}

func (m *mysqlClient) Close() error {
	return m.db.Close()
}

type Tx struct {
	tx *sql.Tx
}

func (t Tx) Commit() error {
	return t.tx.Commit()
}

func (t Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t Tx) Exec(query string, args ...interface{}) (Result, error) {
	return t.tx.Exec(query, args...)
}

func (t Tx) Query(query string, args ...interface{}) (Row, error) {
	return t.tx.Query(query, args...)
}

func (t Tx) QueryRow(query string, args ...interface{}) Row {
	return t.tx.QueryRow(query, args...)
}

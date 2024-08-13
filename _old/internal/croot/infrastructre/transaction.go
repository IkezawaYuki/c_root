package infrastructre

import "database/sql"

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

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
	panic("implement me")
}

func (t Tx) Query(query string, args ...interface{}) (Row, error) {
	//TODO implement me
	panic("implement me")
}

func (t Tx) QueryRow(query string, args ...interface{}) Row {
	//TODO implement me
	panic("implement me")
}

package postgres

import "database/sql"

type DB struct {
	conn *sql.DB
}

func NewDB(conn *sql.DB) *DB {
	return &DB{conn: conn}
}

func (db *DB) Conn() *sql.DB {
	return db.conn
}

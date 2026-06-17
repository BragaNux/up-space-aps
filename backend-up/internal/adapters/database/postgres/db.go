package postgres

import "database/sql"

// DB e um wrapper fino em volta da conexao sql, usado por todos os repositorios postgres
type DB struct {
	conn *sql.DB
}

func NewDB(conn *sql.DB) *DB {
	return &DB{conn: conn}
}

// Conn devolve a conexao sql crua pra rodar queries
func (db *DB) Conn() *sql.DB {
	return db.conn
}

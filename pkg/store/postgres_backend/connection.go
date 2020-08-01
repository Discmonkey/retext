package postgres_backend

import (
	"database/sql"
	"fmt"
	"github.com/discmonkey/retext/pkg/store/credentials"
	_ "github.com/lib/pq"
)

// Connection is a simple utility wrapper around sql.DB
type connection struct {
	db_ *sql.DB
}

func getDb() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		credentials.GetHost(), credentials.GetPort(),
		credentials.GetUser(), credentials.GetPass(), credentials.GetDB())

	return sql.Open("postgres", psqlInfo)
}

func getConnection() (connection, error) {
	c := connection{}

	db, err := getDb()
	if err != nil {
		return c, err
	}

	c.db_ = db

	return c, nil
}

func (c connection) db() *sql.DB {
	return c.db_
}

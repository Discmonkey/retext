package postgres_backend

import (
	"database/sql"
	"fmt"
	"github.com/discmonkey/retext/pkg/store/credentials"
	_ "github.com/lib/pq"
)

type connection = *sql.DB

func GetConnection() (*sql.DB, error) {
	p, err := credentials.GetPort()
	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		credentials.GetHost(), p,
		credentials.GetUser(), credentials.GetPass(), credentials.GetDB())

	return sql.Open("postgres", psqlInfo)
}

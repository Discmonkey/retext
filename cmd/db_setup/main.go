package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/discmonkey/retext/pkg/db/postgres"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
)

func fatalLogIf(err error, message string) {
	if err != nil {
		log.Println(message)
		log.Fatal(err)
	}
}

func initQuery(filepath string) (string, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func getSetupQueryLocation() (string, error) {
	initQueryFilePath := flag.String("init_sql", "", "path to init.sql file")
	flag.Parse()
	if len(*initQueryFilePath) == 0 {
		return "", errors.New("init location not provided")
	}

	return *initQueryFilePath, nil
}

func main() {

	schemaLocation, err := getSetupQueryLocation()
	fatalLogIf(err, "")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		postgres.GetHost(), postgres.GetPort(),
		postgres.GetUser(), postgres.GetPass(), postgres.GetDB())

	db, err := sql.Open("postgres", psqlInfo)
	fatalLogIf(err, "could not open database connection")

	err = db.Ping()
	fatalLogIf(err, "could not ping database")

	tx, err := db.Begin()
	fatalLogIf(err, "could not start transaction")

	query, err := initQuery(schemaLocation)
	fatalLogIf(err, "could not load schema from file")

	_, err = tx.Exec(query)

	if err != nil {
		_ = tx.Rollback()
		fatalLogIf(err, "failed to execute create query")
	}

	err = tx.Commit()
	fatalLogIf(err, "could not commit create schema transaction")
}

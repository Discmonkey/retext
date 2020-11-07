package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/discmonkey/retext/pkg/store/postgres_backend"
	packageVersion "github.com/discmonkey/retext/pkg/version"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func fatalLogIf(err error, message string) {
	if err != nil {
		log.Println(message)
		log.Fatal(err)
	}
}

func sqlFileToString(filepath string) (string, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func getMigrationDir() (string, error) {
	initQueryFilePath := flag.String("migration_dir", "", "path to init.sql file")
	flag.Parse()
	if len(*initQueryFilePath) == 0 {
		return "", errors.New("init location not provided")
	}

	return *initQueryFilePath, nil
}

func getInitSql(migrationDir string) (string, error) {
	fPath := filepath.Join(migrationDir, "init", "init.sql")

	return sqlFileToString(fPath)
}

func getMigrations(migrationDir string) ([]string, error) {
	fileList := make([]string, 0, 0)

	files, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		fileList = append(fileList, f.Name())
	}

	toInt := func(filename string) int {
		i, _ := strconv.ParseInt(strings.TrimSuffix(filename, filepath.Ext(filename)), 10, 32)

		return int(i)
	}

	sort.Slice(fileList, func(i int, j int) bool {
		return toInt(fileList[i]) < toInt(fileList[j])
	})

	sqlFiles := make([]string, len(fileList), len(fileList))

	for i, name := range fileList {
		fString, err := sqlFileToString(path.Join(migrationDir, name))
		if err != nil {
			return nil, err
		}

		sqlFiles[i] = fString
	}

	return sqlFiles, nil
}

func version(tx *sql.Tx) error {
	row := tx.QueryRow("SELECT count(*) from qode.parser WHERE parser.version=$1", packageVersion.Version)
	res := 0

	err := row.Scan(&res)
	if err != nil {
		return err
	}

	if res == 0 {
		_, err = tx.Exec("INSERT INTO qode.parser (version) VALUES ($1)", packageVersion.Version)
		if err != nil {
			return err
		}
	}

	return nil
}
func main() {
	migrationDir, err := getMigrationDir()
	timeoutTries := 10

	fatalLogIf(err, "")

	db, err := postgres_backend.GetConnection()

	fatalLogIf(err, "could not open database connection")

	for ; timeoutTries > 0; timeoutTries-- {
		err = db.Ping()

		if err == nil {
			break
		}
		log.Println("sleeping for 1 second while database starts up")
		time.Sleep(time.Second)
	}

	fatalLogIf(err, "could not ping database")

	tx, err := db.Begin()
	fatalLogIf(err, "could not start transaction")

	query, err := getInitSql(migrationDir)
	fatalLogIf(err, "could not load schema from file")

	queries, err := getMigrations(migrationDir)
	fatalLogIf(err, "could not load migrations")

	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		fatalLogIf(err, "failed to execute create query")
	}

	for _, migration := range queries {
		_, err = tx.Exec(migration)
		if err != nil {
			_ = tx.Rollback()
			fatalLogIf(err, fmt.Sprintf("failed to execute migration: %s", migration))
		}
	}

	err = version(tx)
	if err != nil {
		fatalLogIf(err, "failed to determine correct version")
	}

	err = tx.Commit()
	fatalLogIf(err, "could not commit create schema transaction")
}

package credentials

import (
	"errors"
	"os"
	"strconv"
)

func GetDB() string {
	return "postgres"
}

func GetPass() string {
	return "postgres"
}

func GetUser() string {
	return "postgres"
}

func GetHost() string {
	return "localhost"
}

func GetPort() (int, error) {
	p, isSet := os.LookupEnv("QODE_DB_PORT")
	if isSet {
		i, err := strconv.Atoi(p)
		if err == nil {
			return i, nil
		} else {
			return 0, errors.New("db port flag set to a invalid value")
		}
	}
	return 5432, nil
}

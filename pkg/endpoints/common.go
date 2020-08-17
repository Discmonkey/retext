package endpoints

import (
	"log"
	"net/http"
	"strconv"
)

func HttpNotOk(statusCode int, w http.ResponseWriter, frontendErr string, err error) bool {
	if err != nil {
		http.Error(w, frontendErr, statusCode)
		log.Println(err)
		return true
	} else {
		return false
	}
}

func SliceAtoi(strings []string) ([]int, error) {
	ints := make([]int, len(strings))

	for j, s := range strings {
		i, err := strconv.Atoi(s)

		if err != nil {
			return nil, err
		}

		ints[j] = i
	}

	return ints, nil
}

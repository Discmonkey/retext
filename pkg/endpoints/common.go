package endpoints

import (
	"errors"
	"github.com/discmonkey/retext/pkg/store"
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

const MaxUploadSize = 2 * 1024 * 1024

func GetInt(r *http.Request, key string) (int64, bool) {
	var i int64 = 0
	value, ok := r.URL.Query()[key]

	if !ok {
		return i, ok
	}

	i, err := strconv.ParseInt(value[0], 10, 64)

	if err != nil {
		return i, false
	}

	return i, true
}

func GetIntOk(r *http.Request, w http.ResponseWriter, key, message string) (int64, bool) {
	val, ok := GetInt(r, key)

	if !ok {
		err := errors.New(message)
		HttpNotOk(400, w, err.Error(), err)
	}

	return val, ok
}

func GetStringOk(r *http.Request, w http.ResponseWriter, key, message string) (string, bool) {
	value, ok := r.URL.Query()[key]
	if !ok || len(value) < 1 {
		err := errors.New(message)
		HttpNotOk(400, w, err.Error(), err)
	}

	return value[0], true
}

func ProjectIdOk(r *http.Request, w http.ResponseWriter, errMessage string) (store.ProjectId, bool) {
	projectId, ok := GetInt(r, "project_id")
	if !ok {
		err := errors.New(errMessage)
		HttpNotOk(400, w, err.Error(), err)

		return 0, ok
	}

	return projectId, true
}

func LogIf(err error) {
	if err != nil {
		log.Println(err)
	}
}

func HandleError(message string, w http.ResponseWriter) {
	err := errors.New(message)
	HttpNotOk(400, w, err.Error(), err)
}

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

func GetInt(r *http.Request, key string) (int, bool) {
	var i int64 = 0
	value, ok := r.URL.Query()[key]

	if !ok {
		return int(i), ok
	}

	i, err := strconv.ParseInt(value[0], 10, 64)

	if err != nil {
		return int(i), false
	}

	return int(i), true
}

func GetIntOk(r *http.Request, w http.ResponseWriter, key, message string) (int, bool) {
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
		if !ok {
			err := errors.New("message")
			HttpNotOk(400, w, err.Error(), err)
		}
	}

	return value[0], true
}

func ProjectIdOk(r *http.Request, w http.ResponseWriter, errMessage string) (store.ProjectId, bool) {
	projectId, ok := GetInt(r, "projectId")
	if !ok {
		err := errors.New(errMessage)
		HttpNotOk(400, w, err.Error(), err)

		return 0, ok
	}

	return projectId, true
}

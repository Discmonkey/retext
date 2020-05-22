package endpoints

import (
	"log"
	"net/http"
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

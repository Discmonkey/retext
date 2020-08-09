package code

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"log"
	"net/http"
)

func CreateContainer(store store.CodeStore) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, _ *http.Request) {
		containerID, err := store.CreateContainer()
		if endpoints.HttpNotOk(500, w, "could not create code", err) {
			return
		}

		encoder := json.NewEncoder(w)

		w.Header().Set("Content-Type", "application/json")

		err = encoder.Encode(struct {
			ContainerID int
		}{ContainerID: containerID})

		if err != nil {
			log.Println(err)
		}
	}
}

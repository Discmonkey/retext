package code

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"log"
	"net/http"
)

func CreateContainer(store store.CodeStore) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		projectId, ok := endpoints.ProjectIdOk(r, w, "projectId required to create code")
		if !ok {
			return
		}

		containerId, err := store.CreateContainer(projectId)
		if endpoints.HttpNotOk(500, w, "could not create code", err) {
			return
		}

		encoder := json.NewEncoder(w)

		w.Header().Set("Content-Type", "application/json")

		err = encoder.Encode(struct {
			ContainerId int64
		}{ContainerId: containerId})

		if err != nil {
			log.Println(err)
		}
	}
}

package code

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
)

func Container(codeStore store.CodeStore) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		projectId, ok := endpoints.ProjectIdOk(r, w, "projectId required to create code")
		if !ok {
			return
		}

		containerId, err := codeStore.CreateContainer(projectId)
		if endpoints.HttpNotOk(500, w, "could not create code", err) {
			return
		}

		encoder := json.NewEncoder(w)

		w.Header().Set("Content-Type", "application/json")

		endpoints.LogIf(encoder.Encode(store.CodeContainer{
			Id:         containerId,
			Order:      0,
			Codes:      make([]store.Code, 0, 0),
			Percentage: 0,
		}))
	}
}

package project

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
)

func GetEndpoint(store store.ProjectStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId, ok := endpoints.ProjectIdOk(r, w, "projectId parameter required to load project")
		if !ok {
			return
		}

		project, err := store.GetProject(projectId)
		if endpoints.HttpNotOk(404, w, "could not load project", err) {
			return
		}

		endpoints.LogIf(json.NewEncoder(w).Encode(project))
	}
}

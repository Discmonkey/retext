package file

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
)

type ListResponse = []store.File

func FilesEndpoint(store store.FileStore) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		projectId, ok := endpoints.ProjectIdOk(r, w, "projectId parameter required to list files")
		if !ok {
			return
		}

		w.Header().Set("Content-Type", "application/json")

		files, err := store.GetFiles(projectId)

		if endpoints.HttpNotOk(400, w, "could not load files for project", err) {
			return
		}

		_ = json.NewEncoder(w).Encode(files)
	}

	return t
}

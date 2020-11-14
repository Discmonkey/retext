package file

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"log"
	"net/http"
)

type ListResponse struct {
	Files []store.File `json:"files"`
}

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

		l := ListResponse{}

		files, err := store.GetFiles(projectId)

		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
		} else {
			l.Files = files
			_ = json.NewEncoder(w).Encode(l)
		}
	}

	return t
}
package project

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/store"
	"log"
	"net/http"
)

type ListResponse = []store.Project

func ListEndpoint(store store.ProjectStore) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		projects, err := store.GetProjects()

		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
		} else {
			_ = json.NewEncoder(w).Encode(projects)
		}
	}

	return t
}

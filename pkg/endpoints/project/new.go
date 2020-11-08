package project

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
)

type createRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Month       int    `json:"month"`
	Year        int    `json:"year"`
}

func CreateProject(project store.ProjectStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var req createRequest
		err := json.NewDecoder(r.Body).Decode(&req)

		id, err := project.CreateProject(req.Name, req.Description, req.Month, req.Year)
		if endpoints.HttpNotOk(400, w, "server could not create project", err) {
			return
		}

		project, err := project.GetProject(id)
		if endpoints.HttpNotOk(400, w, "server could not fetch project", err) {
			return
		}

		if endpoints.HttpNotOk(400, w, "error creating project", err) {
			return
		} else {
			_ = json.NewEncoder(w).Encode(project)
		}
	}
}

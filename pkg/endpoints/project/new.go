package project

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
)

func AddCreateProjectEndpoint(project store.ProjectStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		name, ok := endpoints.GetStringOk(r, w, "name", "project name missing from create request")
		if !ok {
			return
		}

		description, ok := endpoints.GetStringOk(r, w, "description", "project description missing from create request")
		if !ok {
			return
		}

		month, ok := endpoints.GetIntOk(r, w, "month", "project month missing form create request")
		if !ok {
			return
		}

		year, ok := endpoints.GetIntOk(r, w, "year", "project year missing from create request")
		if !ok {
			return
		}

		id, err := project.CreateProject(name, description, month, year)
		if endpoints.HttpNotOk(400, w, "error creating project", err) {
			return
		} else {
			_ = json.NewEncoder(w).Encode(struct {
				ProjectId int
			}{ProjectId: id})
		}
	}
}

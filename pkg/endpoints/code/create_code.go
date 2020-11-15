package code

import (
	"encoding/json"
	"errors"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
)

func CreateCode(codeBackend store.CodeStore) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var req store.Code
		err := json.NewDecoder(r.Body).Decode(&req)

		if len(req.Name) == 0 {
			err := errors.New("code parameter required")
			endpoints.HttpNotOk(400, w, err.Error(), err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		req.Id, err = codeBackend.CreateCode(req.Name, req.Container)
		if endpoints.HttpNotOk(400, w, "An error occurred while trying to create the new code.", err) {
			return
		}

		endpoints.LogIf(json.NewEncoder(w).Encode(req))
	}

	return t
}

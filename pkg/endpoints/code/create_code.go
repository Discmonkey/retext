package code

import (
	"encoding/json"
	"errors"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
)

type createRequest struct {
	CodeName    string            `json:"code"`
	ContainerID store.ContainerID `json:"containerID"`
}

func CreateCode(store store.CodeStore) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var req createRequest
		err := json.NewDecoder(r.Body).Decode(&req)

		if len(req.CodeName) == 0 {
			err := errors.New("code parameter required")
			endpoints.HttpNotOk(400, w, err.Error(), err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		codeID, err := store.CreateCode(req.CodeName, req.ContainerID)

		if endpoints.HttpNotOk(400, w, "An error occurred while trying to create the new code.", err) {
			return
		}

		encoder := json.NewEncoder(w)
		newCode, err := store.GetCode(codeID)

		if endpoints.HttpNotOk(400, w, "An error occurred while trying to get the new code.", err) {
			return
		}
		_ = encoder.Encode(newCode)
	}

	return t
}

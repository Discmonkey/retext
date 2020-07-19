package code

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/endpoints"
	"net/http"
)

type associateRequest struct {
	CodeID     db.CodeID         `json:"codeID"`
	DocumentID db.FileID         `json:"key"`
	Text       string            `json:"text"`
	FirstWord  db.WordCoordinate `json:"anchor"`
	LastWord   db.WordCoordinate `json:"last"`
}

func AssociateEndpoint(store db.CodeStore) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var req associateRequest
		err := json.NewDecoder(r.Body).Decode(&req)

		if endpoints.HttpNotOk(400, w, "An error occurred.", err) {
			return
		}

		if req.CodeID == 0 {
			err = errors.New("code parameter required")
			endpoints.HttpNotOk(400, w, err.Error(), err)
			return
		}
		if len(req.DocumentID) == 0 {
			err = errors.New("key parameter required")
			endpoints.HttpNotOk(400, w, err.Error(), err)
			return
		}
		if len(req.Text) == 0 {
			err = errors.New("text parameter required")
			endpoints.HttpNotOk(400, w, err.Error(), err)
			return
		}
		err = store.CodifyText(req.CodeID, req.DocumentID, req.Text, req.FirstWord, req.LastWord)

		if endpoints.HttpNotOk(400, w, "An error occurred while trying to save the selected text.", err) {
			return
		}

		_, _ = fmt.Fprint(w, "")

	}

	return t
}

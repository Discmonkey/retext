package category

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/endpoints"
	"net/http"
)

type associateRequest struct {
	CategoryID db.CategoryID `json:"categoryID"`
	DocumentID db.FileID     `json:"documentID"`
	Text       string        `json:"text"`
}

func AssociateEndpoint(store db.Store) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var req associateRequest
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			endpoints.HttpNotOk(400, w, "", err)
			return
		}

		if len(req.CategoryID) == 0 {
			err = errors.New("category parameter required")
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
		err = store.CategorizeText(req.CategoryID, req.DocumentID, req.Text)

		if endpoints.HttpNotOk(400, w, "", err) {
			return
		}

		_, _ = fmt.Fprint(w, "")

	}

	return t
}

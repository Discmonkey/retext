package code

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
)

type disassociateRequest struct {
	TextIds []int64 `json:"textIds"`
}

func DisassociateText(store store.CodeStore) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var req disassociateRequest
		err := json.NewDecoder(r.Body).Decode(&req)

		if len(req.TextIds) == 0 {
			err := errors.New("textIds parameter required")
			endpoints.HttpNotOk(400, w, err.Error(), err)
			return
		}

		if endpoints.HttpNotOk(400, w, "Invalid TextIds passed", err) {
			return
		}

		err = store.UncodeText(req.TextIds)

		if endpoints.HttpNotOk(400, w, "An error occurred while trying to disassociate text.", err) {
			return
		}

		_, _ = fmt.Fprint(w, "")
	}

	return t
}

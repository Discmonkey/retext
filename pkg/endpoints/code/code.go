package code

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
)

func getCode(w http.ResponseWriter, r *http.Request, codeStore store.CodeStore) {
	id, ok := endpoints.GetInt(r, "id")

	if !ok || id == 0 {
		err := errors.New("id parameter required")
		endpoints.HttpNotOk(http.StatusBadRequest, w, err.Error(), err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	code, err := codeStore.GetCode(id)

	if endpoints.HttpNotOk(400, w, "An error occurred while getting your code. ", err) {
		return
	} else {
		_, _ = fmt.Fprint(w, code)
		_ = json.NewEncoder(w).Encode(code)
	}

}

func postCode(w http.ResponseWriter, r *http.Request, codeStore store.CodeStore) {
	var req store.Code
	err := json.NewDecoder(r.Body).Decode(&req)

	if len(req.Name) == 0 {
		err := errors.New("code parameter required")
		endpoints.HttpNotOk(400, w, err.Error(), err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	req.Id, err = codeStore.CreateCode(req.Name, req.Container)
	if endpoints.HttpNotOk(400, w, "An error occurred while trying to create the new code.", err) {
		return
	}

	endpoints.LogIf(json.NewEncoder(w).Encode(req))
}

func Code(codeBackend store.CodeStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCode(w, r, codeBackend)
		case http.MethodPost:
			postCode(w, r, codeBackend)
		default:
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

}

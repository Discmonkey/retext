package code

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
)

func deleteText(w http.ResponseWriter, r *http.Request, codeStore store.CodeStore) {
	id, ok := endpoints.GetInt(r, "text_id")

	if !ok {
		endpoints.HttpNotOk(400, w, "text_id required to delete text", errors.New("no textId in request"))
	}

	err := codeStore.DeleteText(id)

	if endpoints.HttpNotOk(400, w, "error deleting text", err) {
		return
	}

	_, _ = fmt.Fprint(w, "")
}

func postText(w http.ResponseWriter, r *http.Request, codeStore store.CodeStore) {
	codeId, ok := endpoints.GetInt(r, "code_id")

	if !ok {
		endpoints.HttpNotOk(400, w, "code_id required to delete text", errors.New("no code_id in request"))
	}

	var req store.DocumentText
	err := json.NewDecoder(r.Body).Decode(&req)

	if endpoints.HttpNotOk(400, w, "An error occurred.", err) {
		return
	}

	if req.DocumentId <= 0 {
		endpoints.HandleError("invalid document_id", w)
		return
	}

	if len(req.Text) == 0 {
		endpoints.HandleError("text is required to associate", w)
		return
	}

	if req.FirstWord == nil || req.LastWord == nil {
		endpoints.HandleError("text is required to associate", w)
		return
	}

	textId, err := codeStore.CodifyText(codeId, req.DocumentId, req.Text, *req.FirstWord, *req.LastWord)
	if endpoints.HttpNotOk(400, w, "An error occurred while trying to save the selected text.", err) {
		return
	}

	req.Id = textId
	req.CodeId = codeId

	endpoints.LogIf(json.NewEncoder(w).Encode(req))
}

func Text(codeStore store.CodeStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			deleteText(w, r, codeStore)
		case http.MethodPost:
			postText(w, r, codeStore)
		default:
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

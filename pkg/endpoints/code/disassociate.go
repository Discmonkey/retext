package code

import (
	"errors"
	"fmt"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
)

func DisassociateText(store store.CodeStore) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		strIds, ok := r.URL.Query()["TextIds"]

		if !ok || len(strIds) == 0 {
			err := errors.New("textIds parameter required")
			endpoints.HttpNotOk(400, w, err.Error(), err)
			return
		}

		textIds, err := endpoints.SliceAtoi(strIds)

		if endpoints.HttpNotOk(400, w, "Invalid TextIds passed", err) {
			return
		}

		err = store.UncodeText(textIds)

		if endpoints.HttpNotOk(400, w, "An error occurred while trying to disassociate text.", err) {
			return
		}

		_, _ = fmt.Fprint(w, "")
	}

	return t
}

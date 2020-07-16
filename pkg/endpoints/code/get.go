package code

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/endpoints"
	"net/http"
	"strconv"
)

func GetEndpoint(store db.CodeStore) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		id, err := strconv.Atoi(r.URL.Query().Get("id"))

		if err != nil || id == 0 {
			err := errors.New("id parameter required")
			endpoints.HttpNotOk(http.StatusBadRequest, w, err.Error(), err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		code, err := store.GetCode(id)

		if endpoints.HttpNotOk(400, w, "An error occurred while getting your code. ", err) {
			return
		} else {
			_, _ = fmt.Fprint(w, code)
			_ = json.NewEncoder(w).Encode(code)
		}
	}

	return t
}

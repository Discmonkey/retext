package category

import (
	"errors"
	"fmt"
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/endpoints"
	"net/http"
)

func CreateEndpoint(store db.Store) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		name := r.FormValue("category")

		if len(name) == 0 {
			err := errors.New("category parameter required")
			endpoints.HttpNotOk(400, w, err.Error(), err)
		}

		w.Header().Set("Content-Type", "application/json")

		categoryID, err := store.CreateCategory(name)

		if endpoints.HttpNotOk(400, w, "An error occurred while trying to create the new category.", err) {
			return
		} else {
			_, _ = fmt.Fprint(w, categoryID)
		}
	}

	return t
}

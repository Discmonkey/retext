package category

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/endpoints"
	"net/http"
)

type CategoriesResponse struct {
	Categories []string
}

func ListEndpoint(store db.Store) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		l := CategoriesResponse{}

		categories, err := store.Categories()

		if endpoints.HttpNotOk(400, w, "An error occurred while pulling all categories.", err) {
			return
		} else {
			l.Categories = categories
			_ = json.NewEncoder(w).Encode(l)
		}
	}

	return t
}

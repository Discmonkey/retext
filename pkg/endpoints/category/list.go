package category

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/endpoints"
	"net/http"
)

type CategoriesResponse struct {
	Categories []db.Category `json:"categories"`
}

func ListEndpoint(store db.Store) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		l := CategoriesResponse{}

		categoryIDS, err := store.Categories()

		if endpoints.HttpNotOk(400, w, "An error occurred while pulling all categoryIDS.", err) {
			return
		} else {
			for _, categoryID := range categoryIDS {
				newCat, err := store.GetCategory(categoryID)
				if err != nil {
					endpoints.HttpNotOk(500, w, "Unable to get a category|ID: "+categoryID, err)
					return
				}
				l.Categories = append(l.Categories, newCat)
			}
		}
		_ = json.NewEncoder(w).Encode(l)
	}

	return t
}

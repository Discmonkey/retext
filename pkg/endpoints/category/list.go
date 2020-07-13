package category

import (
	"encoding/json"
	"fmt"
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/endpoints"
	"net/http"
	"sort"
)

type CategoriesResponse = []db.CategoryMain

func ListEndpoint(store db.Store) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		categoryIDs, err := store.Categories()

		if endpoints.HttpNotOk(400, w, "An error occurred while pulling all categoryIDs.", err) {
			return
		}

		l := make(CategoriesResponse, len(categoryIDs))

		for i, categoryID := range categoryIDs {
			main, err := store.GetCategoryMain(categoryID)
			if err != nil {
				endpoints.HttpNotOk(500, w, fmt.Sprintf("Unable to get a category|ID: %d", categoryID), err)
				return
			}

			l[i] = main
		}
		// TODO: user-defined sorting instead of this
		sort.Slice(l, func(i int, j int) bool {
			return l[i].Main < l[j].Main
		})

		_ = json.NewEncoder(w).Encode(l)
	}

	return t
}

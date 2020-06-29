package category

import (
	"encoding/json"
	"fmt"
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/endpoints"
	"net/http"
	"sort"
)

type CategoriesResponse = []db.Category

func ListEndpoint(store db.Store) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		l := CategoriesResponse{}

		categoryIDs, err := store.Categories()

		if endpoints.HttpNotOk(400, w, "An error occurred while pulling all categoryIDs.", err) {
			return
		} else {
			for _, categoryID := range categoryIDs {
				newCat, err := store.GetCategory(categoryID)
				if err != nil {
					endpoints.HttpNotOk(500, w, fmt.Sprintf("Unable to get a category|ID: %d", categoryID), err)
					return
				}
				// subcategories are already included under parent categories
				if !newCat.IsSub {
					l = append(l, newCat)
				}
			}
		}
		sort.Slice(l, func(i int, j int) bool {
			return l[i].ID < l[j].ID
		})
		_ = json.NewEncoder(w).Encode(l)
	}

	return t
}

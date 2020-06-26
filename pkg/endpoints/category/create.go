package category

import (
	"encoding/json"
	"errors"
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/endpoints"
	"net/http"
)

type createRequest struct {
	CategoryName     string        `json:"category"`
	ParentCategoryID db.CategoryID `json:"parentCategoryID"`
}

func CreateEndpoint(store db.Store) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var req createRequest
		err := json.NewDecoder(r.Body).Decode(&req)

		if len(req.CategoryName) == 0 {
			err := errors.New("category parameter required")
			endpoints.HttpNotOk(400, w, err.Error(), err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		categoryID, err := store.CreateCategory(req.CategoryName, req.ParentCategoryID)

		if endpoints.HttpNotOk(400, w, "An error occurred while trying to create the new category.", err) {
			return
		}

		newCat, err := store.GetCategory(categoryID)
		if endpoints.HttpNotOk(400, w, "An error occurred while trying to get the new category.", err) {
			return
		}

		_ = json.NewEncoder(w).Encode(newCat)
	}

	return t
}

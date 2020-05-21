package category

import (
	"fmt"
	"github.com/discmonkey/retext/pkg/db"
	"log"
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
			w.WriteHeader(400)
			log.Println("category parameter required")
		}

		w.Header().Set("Content-Type", "application/json")

		categoryID, err := store.CreateCategory(name)

		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
		} else {
			_, _ = fmt.Fprint(w, categoryID)
		}
	}

	return t
}

package category

import (
	"encoding/json"
	"fmt"
	"github.com/discmonkey/retext/pkg/db"
	"log"
	"net/http"
)

func GetEndpoint(store db.Store) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		id := r.URL.Query().Get("id")

		if len(id) == 0 {
			w.WriteHeader(400)
			log.Println("id parameter required")
		}

		w.Header().Set("Content-Type", "application/json")

		category, err := store.GetCategory(id)

		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
		} else {
			_, _ = fmt.Fprint(w, category)
			_ = json.NewEncoder(w).Encode(category)
		}
	}

	return t
}

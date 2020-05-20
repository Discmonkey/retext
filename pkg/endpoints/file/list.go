package file

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/db"
	"log"
	"net/http"
)

type ListResponse struct {
	Files []string
}

func ListEndpoint(store db.Store) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		l := ListResponse{}

		files, err := store.Files()

		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
		} else {
			l.Files = files
			_ = json.NewEncoder(w).Encode(l)
		}
	}

	return t
}

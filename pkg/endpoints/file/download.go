package file

import (
	"encoding/json"
	"fmt"
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/parser"
	"log"
	"net/http"
)

func DownloadEndpoint(store db.Store) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		keys, ok := r.URL.Query()["key"]

		if !ok || len(keys) != 1 {
			w.WriteHeader(400)
			log.Println("key parameter required for file download request")
		}

		fmt.Println(keys)

		w.Header().Set("Content-Type", "application/json")

		contents, err := store.GetFile(keys[0])

		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
		} else {
			_ = json.NewEncoder(w).Encode(parser.Convert(contents))
		}
	}

	return t
}

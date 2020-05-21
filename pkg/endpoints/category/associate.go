package category

import (
	"fmt"
	"github.com/discmonkey/retext/pkg/db"
	"log"
	"net/http"
)

func AssociateEndpoint(store db.Store) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		categoryID := r.FormValue("category")
		documentID := r.FormValue("key")
		text := r.FormValue("text")

		if len(categoryID) == 0 {
			w.WriteHeader(400)
			log.Println("category parameter required")
			return
		}
		if len(documentID) == 0 {
			w.WriteHeader(400)
			log.Println("key parameter required")
			return
		}
		if len(text) == 0 {
			w.WriteHeader(400)
			log.Println("text parameter required")
			return
		}
		err := store.CategorizeText(categoryID, documentID, text)

		if err != nil {
			fmt.Println(err)
			return
		}

		_, _ = fmt.Fprint(w, "")

	}

	return t
}

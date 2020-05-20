package file

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/db"
	"net/http"
)

package file

import (
"fmt"
"github.com/discmonkey/retext/pkg/db"
"io/ioutil"
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

		json.NewEncoder(w).Encode(user)


	}

	return t
}

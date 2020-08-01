package code

import (
	"encoding/json"
	"fmt"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
	"sort"
)

type CodesResponse = []store.CodeContainer

func ListEndpoint(store store.CodeStore) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		codeIDs, err := store.Codes()

		if endpoints.HttpNotOk(400, w, "An error occurred while pulling all codeIDs.", err) {
			return
		}

		l := make(CodesResponse, len(codeIDs))

		for i, codeID := range codeIDs {
			main, err := store.GetCodeContainer(codeID)
			if err != nil {
				endpoints.HttpNotOk(500, w, fmt.Sprintf("Unable to get a code|ID: %d", codeID), err)
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

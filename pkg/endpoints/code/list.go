package code

import (
	"encoding/json"
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

		// TODO: either add a fileId parameter that will restrict the DocumentWord's returned or add a new endpoint that does. (either way, will still need to return all Codes)
		containers, err := store.GetContainers()

		if endpoints.HttpNotOk(400, w, "could not fetch containers", err) {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		sort.Slice(containers, func(i int, j int) bool {
			return containers[i].Order < containers[j].Order
		})

		_ = json.NewEncoder(w).Encode(containers)
	}

	return t
}

package insight

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
)

func getInsights(w http.ResponseWriter, r *http.Request, insightStore store.InsightStore) {
	projectId, ok := endpoints.ProjectIdOk(r, w, "project Id required to list code containers")
	if !ok {
		return
	}

	insights, err := insightStore.GetInsights(projectId)

	if endpoints.HttpNotOk(400, w, "could not fetch insights", err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(insights)
}

func Insights(insightStore store.InsightStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getInsights(w, r, insightStore)
		default:
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

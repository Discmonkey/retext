package insight

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
)

func createInsight(w http.ResponseWriter, r *http.Request, insightStore store.InsightStore) {
	var req store.Insight
	err := json.NewDecoder(r.Body).Decode(&req)

	if endpoints.HttpNotOk(400, w, "An error occurred", err) {
		return
	}

	if len(req.Value) == 0 {
		endpoints.HandleError("an insight is required to have insight", w)
		return
	}

	projectId, ok := endpoints.ProjectIdOk(r, w, "project Id required to list code containers")
	if !ok {
		return
	}

	insightId, err := insightStore.CreateInsight(projectId, req.Value, req.TextIds)

	if endpoints.HttpNotOk(400, w, "An error occurred while saving your insight", err) {
		return
	}

	req.Id = insightId

	endpoints.LogIf(json.NewEncoder(w).Encode(req))
}

func Insight(insightStore store.InsightStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createInsight(w, r, insightStore)
		default:
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

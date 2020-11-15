package file

import (
	"github.com/discmonkey/retext/pkg/store"
	"net/http"
)

type AddResponse []store.File

func FileEndpoint(fileStore store.FileStore, fileType store.FileType) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			parse(fileStore, fileType, w, r)
		} else if r.Method == http.MethodPost {
			upload(fileStore, fileType, w, r)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

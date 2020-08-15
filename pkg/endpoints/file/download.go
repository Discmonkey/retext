package file

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/parser"
	"github.com/discmonkey/retext/pkg/store"
	"log"
	"net/http"
	"strconv"
)

func fileExtToType(ext string) parser.DocumentType {
	filetype := parser.Text
	if ext == "docx" {
		filetype = parser.DocX
	} else if ext == "xlsx" {
		filetype = parser.Xlsx
	}

	return filetype
}

func DownloadEndpoint(store store.FileStore) func(w http.ResponseWriter, r *http.Request) {
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

		id, err := strconv.ParseInt(keys[0], 10, 64)
		if endpoints.HttpNotOk(400, w, "invalid document id", err) {
			return
		}

		w.Header().Set("Content-Type", "application/json")

		contents, fileSpec, err := store.GetFile(int(id))

		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
		} else {

			filetype := fileExtToType(fileSpec.Ext)

			converted, err := parser.Convert(contents, filetype)

			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				return
			}

			_ = json.NewEncoder(w).Encode(converted)
		}
	}

	return t
}

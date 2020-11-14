package file

import (
	"encoding/json"
	"errors"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/parser"
	"github.com/discmonkey/retext/pkg/store"
	"log"
	"net/http"
)

func fileExtToType(ext string) parser.DocumentType {
	filetype := parser.Text
	if ext == "docx" {
		filetype = parser.DocX
	} else if ext == "xlsx" {
		filetype = parser.Xlsx
	} else if ext == "csv" {
		filetype = parser.Csv
	}

	return filetype
}

func Parse(store store.FileStore, fileType store.FileType, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fileId, ok := endpoints.GetInt(r, "file_id")
	if !ok {
		endpoints.HttpNotOk(400, w, "file_id not present in query", errors.New("missing file_id key"))
	}

	w.Header().Set("Content-Type", "application/json")

	contents, fileSpec, err := store.GetFile(fileId)

	if endpoints.HttpNotOk(400, w, "could not fetch file", err) {
		return
	}

	if fileSpec.Type_ != fileType {
		err = errors.New("file not registered as " + fileType)
	}

	if endpoints.HttpNotOk(400, w, "mismatced file type", err) {
		return
	}

	converted, err := parser.Convert(contents, fileType)

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	_ = json.NewEncoder(w).Encode(converted)

}

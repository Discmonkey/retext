package file

import (
	"encoding/json"
	"fmt"
	"github.com/discmonkey/retext/pkg/db"
	"log"
	"net/http"
)

type Sentence struct {
	Words []string
}

type Paragraph struct {
	Sentences []Sentence
}

type Document struct {
	Tags       []string
	Title      []string
	Paragraphs []Paragraph
}

func convert(unprocessed []byte) Document {

	return Document{}
}

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

		l := DownloadResponse{}

		contents, err := store.GetFile(keys[0])

		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
		} else {
			l.Contents = string(contents)
			_ = json.NewEncoder(w).Encode(l)
		}
	}

	return t
}

package file

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/endpoints"
	"io/ioutil"
	"log"
	"net/http"
)

type AddResponse struct {
	Key  string
	Type string
}

func AddUploadEndpoint(store db.FileStore) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		err := r.ParseMultipartForm(endpoints.MaxUploadSize)
		if err != nil {
			log.Fatal("failed while parsing form")
		}
		form := r.MultipartForm
		fileType := r.FormValue("fileType")
		files := form.File["file"]

		response := make([]AddResponse, 0, len(files))
		for _, handle := range files {
			file, err := handle.Open()
			if err != nil {
				log.Println(err)
				return
			}

			data, err := ioutil.ReadAll(file)
			if err != nil {
				log.Println(err)
				return
			}
			_ = file.Close()

			key, err := store.UploadFile(handle.Filename, data)

			response = append(response, AddResponse{
				Key:  key,
				Type: fileType,
			})
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println(err)
		}
	}

	return t
}

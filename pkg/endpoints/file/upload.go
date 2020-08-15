package file

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/store"
	"io/ioutil"
	"log"
	"net/http"
)

type AddResponse struct {
	File store.File
	Type string
}

func AddUploadEndpoint(store store.FileStore) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		file, handle, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			return
		}

		defer func() {
			err := file.Close()
			if err != nil {
				log.Println(err)
			}
		}()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			return
		}

		uploadedFile, err := store.UploadFile(handle.Filename, data)

		err = json.NewEncoder(w).Encode(AddResponse{
			File: uploadedFile,
			Type: "source",
		})

		if err != nil {
			log.Println(err)
		}

	}

	return t
}

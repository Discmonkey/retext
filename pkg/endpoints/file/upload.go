package file

import (
	"encoding/json"
	"github.com/discmonkey/retext/pkg/endpoints"
	"github.com/discmonkey/retext/pkg/store"
	"io/ioutil"
	"log"
	"net/http"
)

func upload(fileStore store.FileStore, fileType store.FileType, w http.ResponseWriter, r *http.Request) {
	projectId, ok := endpoints.ProjectIdOk(r, w, "projectId required to upload files")
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, endpoints.MaxUploadSize)
	err := r.ParseMultipartForm(endpoints.MaxUploadSize)
	if err != nil {
		log.Println("failed while parsing form")
		return
	}

	form := r.MultipartForm
	files := form.File["files"]

	response := make(AddResponse, 0, len(files))
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

		uploadedFile, err := fileStore.UploadFile(handle.Filename, data, projectId, fileType)

		response = append(response, uploadedFile)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}

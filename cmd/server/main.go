package main

import (
	"github.com/discmonkey/retext/pkg/endpoints/code"
	"github.com/discmonkey/retext/pkg/endpoints/file"
	"github.com/discmonkey/retext/pkg/store/postgres_backend"
	"log"
	"net/http"
	"os"
	"path"
)

func FailIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	fs := http.FileServer(http.Dir("pkg/www/retext/dist"))
	http.Handle("/", fs)

	retextLocation := path.Join(os.TempDir(), "retext")
	log.Printf("file store dir: %s", retextLocation)
	fileBackend, err := postgres_backend.NewFileStore(retextLocation)
	FailIfError(err)

	codeBackend, err := postgres_backend.NewCodeStore()
	FailIfError(err)

	http.HandleFunc("/file/upload", file.AddUploadEndpoint(fileBackend))
	http.HandleFunc("/file/list", file.ListEndpoint(fileBackend))
	http.HandleFunc("/file/load", file.DownloadEndpoint(fileBackend))

	http.HandleFunc("/code/container/create", code.CreateContainer(codeBackend))
	http.HandleFunc("/code/create", code.CreateCode(codeBackend))
	http.HandleFunc("/code/get", code.GetEndpoint(codeBackend))
	http.HandleFunc("/code/list", code.ListEndpoint(codeBackend))
	http.HandleFunc("/code/associate", code.AssociateEndpoint(codeBackend))

	log.Println("Listening on :3000...")
	FailIfError(http.ListenAndServe(":3000", nil))
}

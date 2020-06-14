package main

import (
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/endpoints/category"
	"github.com/discmonkey/retext/pkg/endpoints/file"
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

	backend := &db.FSBackend{}
	retextLocation := path.Join(os.TempDir(), "retext")
	FailIfError(backend.Init(retextLocation))

	http.HandleFunc("/file/upload", file.AddUploadEndpoint(backend))

	http.HandleFunc("/file/list", file.ListEndpoint(backend))

	http.HandleFunc("/file/load", file.DownloadEndpoint(backend))

	http.HandleFunc("/category/create", category.CreateEndpoint(backend))

	http.HandleFunc("/category/get", category.GetEndpoint(backend))

	http.HandleFunc("/category/list", category.ListEndpoint(backend))

	http.HandleFunc("/category/associate", category.AssociateEndpoint(backend))

	log.Println("Listening on :3000...")
	FailIfError(http.ListenAndServe(":3000", nil))
}

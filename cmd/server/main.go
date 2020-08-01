package main

import (
	"github.com/discmonkey/retext/pkg/endpoints/code"
	"github.com/discmonkey/retext/pkg/endpoints/file"
	"github.com/discmonkey/retext/pkg/store/file_backend"
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

	fileBackend := &file_backend.DevFileBackend{}
	FailIfError(fileBackend.Init(retextLocation))

	codeBackend := &file_backend.DevCodeBackend{}
	FailIfError(codeBackend.Init(retextLocation))

	http.HandleFunc("/file/upload", file.AddUploadEndpoint(fileBackend))

	http.HandleFunc("/file/list", file.ListEndpoint(fileBackend))

	http.HandleFunc("/file/load", file.DownloadEndpoint(fileBackend))

	http.HandleFunc("/code/create", func(writer http.ResponseWriter, request *http.Request) {
		code.CreateEndpoint(codeBackend)(writer, request)
	})
	http.HandleFunc("/code/get", func(writer http.ResponseWriter, request *http.Request) {
		code.GetEndpoint(codeBackend)(writer, request)
	})
	http.HandleFunc("/code/list", func(writer http.ResponseWriter, request *http.Request) {
		code.ListEndpoint(codeBackend)(writer, request)
	})
	http.HandleFunc("/code/associate", func(writer http.ResponseWriter, request *http.Request) {
		code.AssociateEndpoint(codeBackend)(writer, request)
	})

	log.Println("Listening on :3000...")
	FailIfError(http.ListenAndServe(":3000", nil))
}

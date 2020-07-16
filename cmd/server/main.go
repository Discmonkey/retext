package main

import (
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/endpoints/code"
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

	retextLocation := path.Join(os.TempDir(), "retext")

	fileBackend := &db.DevFileBackend{}
	FailIfError(fileBackend.Init(retextLocation))

	codeBackend := &db.DevCodeBackend{}
	FailIfError(codeBackend.Init(retextLocation))

	http.HandleFunc("/file/upload",
		func(writer http.ResponseWriter, request *http.Request) {
			enableCors(&writer)
			file.AddUploadEndpoint(fileBackend)(writer, request)
		})
	http.HandleFunc("/file/list", func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		file.ListEndpoint(fileBackend)(writer, request)
	})
	http.HandleFunc("/file/load", func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		file.DownloadEndpoint(fileBackend)(writer, request)
	})

	http.HandleFunc("/code/create", func(writer http.ResponseWriter, request *http.Request) {
		code.CreateEndpoint(codeBackend)(writer, request)
	})
	http.HandleFunc("/code/get", func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		code.GetEndpoint(codeBackend)(writer, request)
	})
	http.HandleFunc("/code/list", func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		code.ListEndpoint(codeBackend)(writer, request)
	})
	http.HandleFunc("/code/associate", func(writer http.ResponseWriter, request *http.Request) {
		code.AssociateEndpoint(codeBackend)(writer, request)
	})

	log.Println("Listening on :3000...")
	FailIfError(http.ListenAndServe(":3000", nil))
}

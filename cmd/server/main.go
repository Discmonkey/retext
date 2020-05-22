package main

import (
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/endpoints/category"
	"github.com/discmonkey/retext/pkg/endpoints/file"
	"log"
	"net/http"
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
	FailIfError(backend.Init("/tmp/uploadLocation"))

	http.HandleFunc("/file/upload",
		func(writer http.ResponseWriter, request *http.Request) {
			enableCors(&writer)
			file.AddUploadEndpoint(backend)(writer, request)
		})
	http.HandleFunc("/file/list", func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		file.ListEndpoint(backend)(writer, request)
	})
	http.HandleFunc("/file/load", func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		file.DownloadEndpoint(backend)(writer, request)
	})

	http.HandleFunc("/category/create", func(writer http.ResponseWriter, request *http.Request) {
		category.CreateEndpoint(backend)(writer, request)
	})
	http.HandleFunc("/category/get", func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		category.GetEndpoint(backend)(writer, request)
	})
	http.HandleFunc("/category/list", func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		category.ListEndpoint(backend)(writer, request)
	})
	http.HandleFunc("/category/associate", func(writer http.ResponseWriter, request *http.Request) {
		category.AssociateEndpoint(backend)(writer, request)
	})

	log.Println("Listening on :3000...")
	FailIfError(http.ListenAndServe(":3000", nil))
}

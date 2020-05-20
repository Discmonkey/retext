package main

import (
	"github.com/discmonkey/retext/pkg/db"
	"github.com/discmonkey/retext/pkg/endpoints/file"
	"log"
	"net/http"
)

func FailIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	fs := http.FileServer(http.Dir("pkg/www/retext/dist"))
	http.Handle("/", fs)

	backend := &db.FSBackend{}
	FailIfError(backend.Init("/tmp/uploadLocation"))

	http.HandleFunc("/file/upload", file.AddUploadEndpoint(backend))

	log.Println("Listening on :3000...")
	FailIfError(http.ListenAndServe(":3000", nil))
}

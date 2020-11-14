package main

import (
	"flag"
	"github.com/discmonkey/retext/pkg/endpoints/code"
	"github.com/discmonkey/retext/pkg/endpoints/file"
	"github.com/discmonkey/retext/pkg/endpoints/project"
	"github.com/discmonkey/retext/pkg/store"
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

func getTempFileDir() string {
	retextLocation := path.Join(os.TempDir(), "retext")

	if _, err := os.Stat(retextLocation); os.IsNotExist(err) {
		_ = os.Mkdir(retextLocation, 0775)
	}

	return retextLocation
}

func getFileSaveDir() string {
	location := flag.String("file_dir", "", "directory where files are stored")

	flag.Parse()

	if len(*location) == 0 {
		return getTempFileDir()
	} else {
		return *location
	}
}

func main() {
	retextLocation := getFileSaveDir()
	log.Printf("file store dir: %s", retextLocation)

	connection, err := postgres_backend.GetConnection()
	FailIfError(err)

	fs := http.FileServer(http.Dir("pkg/www/retext/dist"))
	http.Handle("/", fs)

	fileBackend := postgres_backend.NewFileStore(retextLocation, connection)
	codeBackend := postgres_backend.NewCodeStore(connection)
	projectBackend := postgres_backend.NewProjectStore(connection)

	http.HandleFunc("/files", file.FilesEndpoint(fileBackend))
	http.HandleFunc("/demo", file.FileEndpoint(fileBackend, store.DemoFile))
	http.HandleFunc("/source", file.FileEndpoint(fileBackend, store.SourceFile))

	http.HandleFunc("/code/container/create", code.CreateContainer(codeBackend))
	http.HandleFunc("/code/create", code.CreateCode(codeBackend))
	http.HandleFunc("/code/get", code.GetEndpoint(codeBackend))
	http.HandleFunc("/code/list", code.ListEndpoint(codeBackend))
	http.HandleFunc("/code/associate", code.AssociateEndpoint(codeBackend))
	http.HandleFunc("/code/disassociate", code.DisassociateText(codeBackend))

	http.HandleFunc("/project", project.GetEndpoint(projectBackend))
	http.HandleFunc("/project/create", project.CreateProject(projectBackend))
	http.HandleFunc("/project/list", project.ListEndpoint(projectBackend))

	log.Println("Listening on :3000...")
	FailIfError(http.ListenAndServe(":3000", nil))
}

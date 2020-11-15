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

// from NotFoundRedirectRespWr to wrapHandler mostly from https://stackoverflow.com/questions/47285119/
type NotFoundRedirectRespWr struct {
	http.ResponseWriter // We embed http.ResponseWriter
	status              int
}

func (w *NotFoundRedirectRespWr) WriteHeader(status int) {
	w.status = status // Store the status for our own use
	if status != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(status)
	}
}

const staticDir = "pkg/www/retext/dist"

func (w *NotFoundRedirectRespWr) Write(p []byte) (int, error) {
	if w.status != http.StatusNotFound {
		return w.ResponseWriter.Write(p)
	}
	return len(p), nil // Lie that we successfully written it
}

//wrapHandler we have a single-page application; chances are good that the url doesn't
// correspond to an actual page. Instead of 404, just serve the index page and hope for the best.
// (assume the front-end will 404 as necessary)
func wrapHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nfrw := &NotFoundRedirectRespWr{ResponseWriter: w}
		h.ServeHTTP(nfrw, r)
		if nfrw.status == 404 {
			// the call to h.serveHTTP above sets the C-t to text/plain. Overwrite to appropriate type.
			w.Header().Add("Content-type", "text/html")

			f, err := os.Open(staticDir + "/index.html")
			if err != nil {
				// index file not found? just give up cause something is terribly wrong
				log.Fatalf("error opening index.html: %s", err)
			}
			fs, err := f.Stat()
			if err != nil {
				// no stats for index file? just give up cause something is terribly wrong
				log.Fatalf("error getting stats for index.html: %s", err)
			}

			http.ServeContent(w, r, "index.html", fs.ModTime(), f)
		}
	}
}

func main() {
	retextLocation := getFileSaveDir()
	log.Printf("file store dir: %s", retextLocation)

	connection, err := postgres_backend.GetConnection()
	FailIfError(err)

	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", wrapHandler(fs))

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

package file

import (
	"fmt"
	"github.com/discmonkey/retext/pkg/db"
	"io/ioutil"
	"net/http"
)

func AddUploadEndpoint(store db.Store) func(w http.ResponseWriter, r *http.Request) {
	t := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		file, handle, err := r.FormFile("file")
		if err != nil {
			fmt.Println(w, "%v", err)
			return
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			return
		}

		key, err := store.UploadFile(handle.Filename, data)

		fmt.Fprint(w, key)

	}

	return t
}

package postgres_backend

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/discmonkey/retext/pkg/store"
	"github.com/discmonkey/retext/pkg/store/file_backend"
	"io/ioutil"
	"path"
	"strings"
)

type FileStore struct {
	fileSys  *file_backend.DevFileBackend
	db       connection
	writeDir string
}

func (c FileStore) UploadFile(filename string, contents []byte) (store.FileID, error) {

	// generate has for file contents
	hash := hashContents(contents)

	exists, err := checkExists(c.db, hash)
	if err != nil {
		return "", err
	}

	var location string

	if exists {
		location, err = getLocationFromHash(c.db, hash)
		if err != nil {
			return "", err
		}
	} else {
		_, err = c.fileSys.UploadFile(filename, contents)
		if err != nil {
			return "", err
		}

		location = path.Join(c.writeDir, filename)
	}

	return logFileToDb(c.db, filename, location, hash)
}

func (c FileStore) GetFile(id store.FileID) ([]byte, error) {
	location, err := getLocationFromID(c.db, id)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(location)

}

func (c FileStore) Files() ([]store.File, error) {
	return listFiles(c.db)
}

func NewStore(writeLocation string) (*FileStore, error) {
	fStore := file_backend.DevFileBackend{}

	err := fStore.Init(writeLocation)
	if err != nil {
		return nil, err
	}

	c := FileStore{
		writeDir: writeLocation,
		fileSys:  &fStore,
	}

	con, err := getConnection()
	if err != nil {
		return nil, err
	}

	c.db = con

	return &c, nil
}

func hashContents(contents []byte) string {
	h := sha256.New()
	h.Write(contents)

	return hex.EncodeToString(h.Sum(nil))
}

func checkExists(con connection, fileContentsHash string) (bool, error) {
	query := `SELECT count(*) FROM qode.file WHERE file_hash = $1`

	row := con.db().QueryRow(query, fileContentsHash)

	var num int
	err := row.Scan(&num)

	if err != nil {
		return false, err
	}

	return num > 0, nil
}

func getLocationFromHash(con connection, fileContentsHash string) (string, error) {

	query := `SELECT location FROM qode.file WHERE file_hash = ?`

	row := con.db().QueryRow(query, fileContentsHash)

	var id int
	err := row.Scan(&id)

	return fmt.Sprintf("%d", id), err
}

func getLocationFromID(con connection, id store.FileID) (string, error) {
	query := `SELECT location FROM qode.file WHERE id = $1`

	row := con.db().QueryRow(query, id)

	var location string
	err := row.Scan(&location)

	return fmt.Sprintf("%s", id), err
}

func logFileToDb(con connection, filename, location, hash string) (store.FileID, error) {
	insert := `
		INSERT INTO qode.file (name, uploaded, location, file_hash) 
		VALUES ($1, NOW(), $2, $3)
		RETURNING id 
		`

	var id int

	row := con.db().QueryRow(insert, filename, location, hash)

	err := row.Scan(&id)

	return fmt.Sprintf("%d", id), err
}

func listFiles(con connection) ([]store.File, error) {
	query := `
		SELECT id, name as string FROM qode.file 
	`

	res := make([]store.File, 0, 0)

	rows, err := con.db().Query(query)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	var id int
	var filename string

	for rows.Next() {
		err = rows.Scan(&id, &filename)
		if err != nil {
			return nil, err
		}

		f := store.File{}
		f.ID = fmt.Sprintf("%d", id)

		f.Type = store.SourceFile
		if strings.HasSuffix(filename, "xlsx") {
			f.Type = store.DemoFile
		}

		res = append(res, f)
	}

	return res, nil
}

var _ store.FileStore = FileStore{}

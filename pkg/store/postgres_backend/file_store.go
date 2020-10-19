package postgres_backend

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/discmonkey/retext/pkg/store"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

type FileSys interface {
	Store(filepath string, contents []byte) error
	Fetch(filepath string) ([]byte, error)
}

type DefaultFileSys struct {
}

func (d DefaultFileSys) Store(filepath string, contents []byte) error {
	return ioutil.WriteFile(filepath, contents, 0666)
}

func (d DefaultFileSys) Fetch(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

var _ FileSys = DefaultFileSys{}

type FileStore struct {
	fileSys  FileSys
	db       *sql.DB
	writeDir string
}

func (c FileStore) UploadFile(filename string, contents []byte, projectId store.ProjectId) (store.File, error) {

	// generate has for file contents
	hash := hashContents(contents)

	exists, err := checkExists(c.db, hash)
	if err != nil {
		return store.File{}, err
	}

	var location string

	if exists {
		location, err = getLocationFromHash(c.db, hash)
		if err != nil {
			return store.File{}, err
		}
	} else {
		location = path.Join(c.writeDir, filename)
		err = c.fileSys.Store(location, contents)
		if err != nil {
			return store.File{}, err
		}
	}

	id, err := logFileToDb(c.db, filename, location, hash, projectId)
	if err != nil {
		return store.File{}, err
	}

	name, ext, err := getNameAndExtension(filename)
	if err != nil {
		return store.File{}, err
	}

	return store.File{
		Id:   id,
		Type: assignTypeFromExtension(ext),
		Name: name,
	}, nil
}

func (c FileStore) GetFile(id store.FileId) ([]byte, store.File, error) {
	query := `SELECT id, name, location as string FROM qode.file WHERE id = $1`
	var location string
	row := c.db.QueryRow(query, id)

	file, err := parseFileRow(row, &location)
	if err != nil {
		return nil, store.File{}, err
	}

	contents, err := c.fileSys.Fetch(location)

	return contents, file, err
}

func (c FileStore) GetFiles(id store.ProjectId) ([]store.File, error) {
	return listFiles(c.db, id)
}

func NewFileStore(writeLocation string) (*FileStore, error) {
	c := FileStore{
		writeDir: writeLocation,
		fileSys:  DefaultFileSys{},
	}

	con, err := GetConnection()
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

	row := con.QueryRow(query, fileContentsHash)

	var num int
	err := row.Scan(&num)

	if err != nil {
		return false, err
	}

	return num > 0, nil
}

func getLocationFromHash(con connection, fileContentsHash string) (string, error) {

	query := `SELECT location FROM qode.file WHERE file_hash = $1`

	row := con.QueryRow(query, fileContentsHash)

	var location string
	err := row.Scan(&location)

	return location, err
}

func getLocationFromId(con connection, id store.FileId) (string, error) {
	query := `SELECT location FROM qode.file WHERE id = $1`

	row := con.QueryRow(query, id)

	var location string
	err := row.Scan(&location)

	return fmt.Sprintf("%s", location), err
}

func logFileToDb(con connection, filename, location, hash string, project store.ProjectId) (store.FileId, error) {
	insert := `
		INSERT INTO qode.file (name, uploaded, location, file_hash, project_id) 
		VALUES ($1, NOW(), $2, $3, $4)
		RETURNING id 
		`

	var id int

	row := con.QueryRow(insert, filename, location, hash, project)

	err := row.Scan(&id)

	return id, err
}

func getNameAndExtension(filename string) (string, string, error) {
	nameAndExtension := strings.Split(filename, ".")

	if len(nameAndExtension) != 2 {
		return "", "", errors.New(fmt.Sprintf("file %s cannot be parsed into name and extension", filename))
	}

	return nameAndExtension[0], nameAndExtension[1], nil
}

func assignTypeFromExtension(ext string) store.FileType {
	type_ := store.SourceFile
	if ext == "xlsx" || ext == "csv" {
		type_ = store.DemoFile
	}

	return type_
}

func listFiles(con connection, project store.ProjectId) ([]store.File, error) {
	query := `
		SELECT id, name, location as string FROM qode.file
		WHERE project_id = $1
	`

	res := make([]store.File, 0, 0)

	rows, err := con.Query(query, project)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()
	// not used for now
	var location string
	for rows.Next() {

		file, err := parseFileRow(rows, &location)

		if err != nil {
			log.Println(err)
			continue
		}

		res = append(res, file)
	}

	return res, nil
}

type rowLike interface {
	Scan(dest ...interface{}) error
}

func parseFileRow(row rowLike, location *string) (store.File, error) {
	var filename string
	var id int

	err := row.Scan(&id, &filename, location)

	if err != nil {
		return store.File{}, err
	}

	name, ext, err := getNameAndExtension(filename)

	if err != nil {
		return store.File{}, err
	}

	f := store.File{}
	f.Id = id
	f.Name = name
	f.Type = assignTypeFromExtension(ext)
	f.Ext = ext

	return f, nil
}

var _ store.FileStore = FileStore{}

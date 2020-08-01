package postgres_backend

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/discmonkey/retext/pkg/store"
	"github.com/discmonkey/retext/pkg/store/file_backend"
)

type CodeStore struct {
	fileSys *file_backend.DevFileBackend
	db      connection
}

func NewStore(writeLocation string) (*CodeStore, error) {
	fStore := file_backend.DevFileBackend{}

	err := fStore.Init(writeLocation)
	if err != nil {
		return nil, err
	}

	c := CodeStore{}

	con, err := getConnection()
	if err != nil {
		return nil, err
	}

	c.db = con

	return &c, nil
}

func hash(filename string) string {
	h := sha1.New()
	h.Write([]byte(filename))
	return hex.EncodeToString(h.Sum(nil))
}

func (c CodeStore) UploadFile(filename string, contents []byte) (store.FileID, error) {
	filenameHash := hash(filename)

	checkExistsQuery := ``
	insertQuery := ``

	path, err := c.fileSys.UploadFile(filenameHash, contents)
	if err != nil {
		return "", err
	}

}

func (c CodeStore) GetFile(id store.FileID) ([]byte, error) {
	panic("implement me")
}

func (c CodeStore) Files() ([]store.File, error) {
	panic("implement me")
}

var _ store.FileStore = CodeStore{}

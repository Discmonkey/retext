package db

type ID = string

type Store interface {
	UploadFile(filename string, contents []byte) (ID, error)
	GetFile(id ID) ([]byte, error)
	Files() ([]ID, error)
}

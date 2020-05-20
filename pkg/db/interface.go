package db

type FileID = string
type CategoryID = string
type Category struct {
	Id   CategoryID
	Name string
}

type Store interface {
	UploadFile(filename string, contents []byte) (FileID, error)
	GetFile(id FileID) ([]byte, error)
	Files() ([]FileID, error)
	CreateCategory(name string) (Category, error)
	CategorizeText(categoryID CategoryID, documentID FileID, text string) error
	GetCategory(categoryID CategoryID) (Category, error)
	Categories() ([]CategoryID, error)
}

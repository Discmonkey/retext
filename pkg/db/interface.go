package db

type ID = string
type Category struct {
	id   ID
	name string
}

type Store interface {
	UploadFile(filename string, contents []byte) (ID, error)
	GetFile(id ID) ([]byte, error)
	Files() ([]ID, error)
	CreateCategory(name string) (Category, error)
	CategorizeText(categoryId ID, documentTexts []map[ID]string) error
	GetCategory(ID) (Category, error)
	Categories() ([]ID, error)
}

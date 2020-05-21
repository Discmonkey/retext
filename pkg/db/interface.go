package db

type FileID = string

type CategoryID = string
type DocumentText struct {
	DocumentID FileID
	Text       string
}
type Category struct {
	ID    CategoryID     `json:"id"`
	Name  string         `json:"name"`
	Texts []DocumentText `json:"texts"`
}
type Categories struct {
	Categories map[string]Category `json:"categories"`
}

type Store interface {
	UploadFile(filename string, contents []byte) (FileID, error)
	GetFile(id FileID) ([]byte, error)
	Files() ([]FileID, error)
	CreateCategory(name string) (CategoryID, error)
	CategorizeText(categoryID CategoryID, documentID FileID, text string) error
	GetCategory(categoryID CategoryID) (Category, error)
	Categories() ([]CategoryID, error)
}
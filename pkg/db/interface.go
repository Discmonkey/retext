package db

type FileID = string
type CategoryID = int

type DocumentText struct {
	DocumentID FileID `json:"documentID"`
	Text       string `json:"text"`
}

type Category struct {
	ID    CategoryID     `json:"id"`
	Name  string         `json:"name"`
	Texts []DocumentText `json:"texts"`
}

type Categories = map[CategoryID]Category

type Store interface {
	UploadFile(filename string, contents []byte) (FileID, error)
	//DeleteFile(id FileID) error
	GetFile(id FileID) ([]byte, error)
	Files() ([]FileID, error)
	CreateCategory(name string) (CategoryID, error)
	CategorizeText(categoryID CategoryID, documentID FileID, text string) error
	GetCategory(categoryID CategoryID) (Category, error)
	Categories() ([]CategoryID, error)
}

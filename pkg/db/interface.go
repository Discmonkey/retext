package db

type FileID = string

type CategoryID = int

type WordCoordinate struct {
	Paragraph int `json:"paragraph"`
	Sentence  int `json:"sentence"`
	Word      int `json:"word"`
}

type DocumentText struct {
	DocumentID FileID `json:"documentID"`
	// TODO: don't store text (since the WordCoordinate's are already stored)
	Text      string         `json:"text"`
	FirstWord WordCoordinate `json:"anchor"`
	LastWord  WordCoordinate `json:"last"`
}

type Category struct {
	ID    CategoryID     `json:"id"`
	Name  string         `json:"name"`
	Texts []DocumentText `json:"texts"`
}
type CategoryMain struct {
	Main       CategoryID `json:"main"`
	Categories []Category `json:"subcategories"`
}

type CategoryParentIDMap = map[CategoryID][]CategoryID
type CategoryMap = map[CategoryID]Category

type Store interface {
	UploadFile(filename string, contents []byte) (FileID, error)
	GetFile(id FileID) ([]byte, error)
	Files() ([]FileID, error)
	CreateCategory(name string, ParentCategoryID CategoryID) (CategoryID, error)
	CategorizeText(categoryID CategoryID, documentID FileID, text string, firstWord WordCoordinate, lastWord WordCoordinate) error
	GetCategory(categoryID CategoryID) (Category, error)
	GetCategoryMain(categoryID CategoryID) (CategoryMain, error)
	Categories() ([]CategoryID, error)
}

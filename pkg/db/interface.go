package db

type FileID = string
type FileType = string

const (
	SourceFile FileType = "SourceFile"
	DemoFile   FileType = "DemoFile"
)

type File struct {
	Id   FileID
	Type FileType
}

type CategoryID = int

type WordCoordinate struct {
	Paragraph int `json:"paragraph"`
	Sentence  int `json:"sentence"`
	Word      int `json:"word"`
}

type DocumentText struct {
	DocumentID FileID         `json:"documentID"`
	Text       string         `json:"text"`
	FirstWord  WordCoordinate `json:"anchor"`
	LastWord   WordCoordinate `json:"last"`
}

type Category struct {
	ID            CategoryID     `json:"id"`
	Name          string         `json:"name"`
	Texts         []DocumentText `json:"texts"`
	IsSub         bool           `json:"isSub"`
	Subcategories []*Category
}

type CategoryContainer struct {
	main     Category
	children []Category
}

type Subcategories = []*Category
type CategoryParentMap struct {
	parentID CategoryID
	childID  CategoryID
}
type Categories = map[CategoryID]*Category

type Store interface {
	UploadFile(filename string, contents []byte) (FileID, error)
	GetFile(id FileID) ([]byte, error)
	Files() ([]File, error)
	CreateCategory(name string, ParentCategoryID CategoryID) (CategoryID, error)
	CategorizeText(categoryID CategoryID, documentID FileID, text string, firstWord WordCoordinate, lastWord WordCoordinate) error
	GetCategory(categoryID CategoryID) (Category, error)
	Categories() ([]CategoryID, error)
}

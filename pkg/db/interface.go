package db

type FileID = string

type CodeID = int

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

type Code struct {
	ID    CodeID         `json:"id"`
	Name  string         `json:"name"`
	Texts []DocumentText `json:"texts"`
}
type CodeContainer struct {
	Main  CodeID `json:"main"`
	Codes []Code `json:"subcodes"`
}

type CodeParentIDMap = map[CodeID][]CodeID
type CodeMap = map[CodeID]Code

type FileStore interface {
	UploadFile(filename string, contents []byte) (FileID, error)
	GetFile(id FileID) ([]byte, error)
	Files() ([]FileID, error)
}

type CodeStore interface {
	CreateCode(name string, ParentCodeID CodeID) (CodeID, error)
	CategorizeText(codeID CodeID, documentID FileID, text string, firstWord WordCoordinate, lastWord WordCoordinate) error
	GetCode(codeID CodeID) (Code, error)
	GetCodeContainer(codeID CodeID) (CodeContainer, error)
	Codes() ([]CodeID, error)
}

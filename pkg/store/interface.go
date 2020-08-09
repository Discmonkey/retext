package store

type FileID = int
type FileType = string

type CodeID = int
type ContainerID = int

const (
	SourceFile FileType = "SourceFile"
	DemoFile   FileType = "DemoFile"
)

type File struct {
	ID   FileID
	Type FileType
	Name string
}

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

type Code struct {
	ID        CodeID         `json:"id"`
	Name      string         `json:"name"`
	Texts     []DocumentText `json:"texts"`
	Container ContainerID    `json:"container"`
}

type CodeContainer struct {
	Main  CodeID `json:"main"`
	Order int
	Codes []Code `json:"subcodes"`
}

type CodeParentIDMap = map[CodeID][]CodeID
type CodeMap = map[CodeID]Code

type FileStore interface {
	UploadFile(filename string, contents []byte) (File, error)
	GetFile(id FileID) ([]byte, error)
	Files() ([]File, error)
}

type CodeStore interface {
	CreateContainer() (ContainerID, error)
	CreateCode(name string, containerID ContainerID) (CodeID, error)
	CodifyText(codeID CodeID, documentID FileID, text string, firstWord WordCoordinate, lastWord WordCoordinate) error
	GetCode(codeID CodeID) (Code, error)
	GetContainer(codeID ContainerID) (CodeContainer, error)
	GetContainers() ([]CodeContainer, error)
}

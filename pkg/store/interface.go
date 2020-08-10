package store

import "time"

type ProjectId = int

type FileId = int
type FileType = string

type CodeId = int
type ContainerId = int

const (
	SourceFile FileType = "SourceFile"
	DemoFile   FileType = "DemoFile"
)

type Project struct {
	Id          ProjectId
	TimeTag     time.Time
	Name        string
	Description string
}

type File struct {
	Id   FileId
	Type FileType
	Name string
	Ext  string
}

type WordCoordinate struct {
	Paragraph int `json:"paragraph"`
	Sentence  int `json:"sentence"`
	Word      int `json:"word"`
}

type DocumentText struct {
	DocumentId FileId         `json:"documentId"`
	Text       string         `json:"text"`
	FirstWord  WordCoordinate `json:"anchor"`
	LastWord   WordCoordinate `json:"last"`
}

type Code struct {
	Id        CodeId         `json:"id"`
	Name      string         `json:"name"`
	Texts     []DocumentText `json:"texts"`
	Container ContainerId    `json:"containerId"`
}

type CodeContainer struct {
	Id    ContainerId `json:"containerId"`
	Order int
	Codes []Code `json:"subcodes"`
}

type CodeParentIdMap = map[CodeId][]CodeId
type CodeMap = map[CodeId]Code

type FileStore interface {
	UploadFile(filename string, contents []byte) (File, error)
	GetFile(id FileId) ([]byte, File, error)
	Files() ([]File, error)
}

type CodeStore interface {
	CreateContainer(id ProjectId) (ContainerId, error)
	CreateCode(name string, containerId ContainerId) (CodeId, error)
	CodifyText(codeId CodeId, documentId FileId, text string, firstWord WordCoordinate, lastWord WordCoordinate) error
	GetCode(codeId CodeId) (Code, error)
	GetContainer(codeId ContainerId) (CodeContainer, error)
	GetContainers(id ProjectId) ([]CodeContainer, error)
}

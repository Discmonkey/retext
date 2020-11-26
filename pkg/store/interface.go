package store

import (
	swagger "github.com/discmonkey/retext/pkg/swagger"
)

type ProjectId = int64

type FileId = int64
type FileType = string

type CodeId = int64
type ContainerId = int64
type TextId = int64
type InsightId = int64

const (
	SourceFile FileType = "KSOURCE"
	DemoFile   FileType = "KDEMO"
)

type Project = swagger.Project
type File = swagger.File
type WordCoordinate = swagger.WordCoordinate
type DocumentText = swagger.DocumentText
type Code = swagger.Code
type CodeContainer = swagger.CodeContainer
type Insight = swagger.Insight

type CodeParentIdMap = map[CodeId][]CodeId
type CodeMap = map[CodeId]Code

type ProjectStore interface {
	CreateProject(name, description string, month, year int) (ProjectId, error)
	GetProject(id ProjectId) (Project, error)
	GetProjects() ([]Project, error)
}

type FileStore interface {
	UploadFile(filename string, contents []byte, id ProjectId, fileType FileType) (File, error)
	GetFile(id FileId) ([]byte, File, error)
	GetFiles(id ProjectId) ([]File, error)
}

type CodeStore interface {
	CreateContainer(id ProjectId) (ContainerId, error)
	CreateCode(name string, containerId ContainerId) (CodeId, error)
	CodifyText(codeId CodeId, documentId FileId, text string, firstWord WordCoordinate, lastWord WordCoordinate) (TextId, error)
	DeleteText(textId TextId) error
	GetCode(codeId CodeId) (Code, error)
	GetContainer(codeId ContainerId) (CodeContainer, error)
	GetContainers(id ProjectId) ([]CodeContainer, error)
}

type InsightStore interface {
	CreateInsight(projectId ProjectId, insightText string, textIds []TextId) (InsightId, error)
	GetInsights(projectId ProjectId) ([]Insight, error)
}

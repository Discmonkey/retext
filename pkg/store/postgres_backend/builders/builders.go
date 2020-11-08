package builders

import (
	"database/sql"
	"github.com/discmonkey/retext/pkg/store"
	"math"
)

type CodeRow struct {
	Name sql.NullString

	P1 sql.NullInt32
	S1 sql.NullInt32
	W1 sql.NullInt32

	P2 sql.NullInt32
	S2 sql.NullInt32
	W2 sql.NullInt32

	Text     sql.NullString
	SourceId sql.NullInt32
	TextId   sql.NullInt32
}

type ContainerRow struct {
	CodeRow          CodeRow
	CodeDisplayOrder int
	CodeId           int
}

type ContainerListRow struct {
	ContainerId    int
	ContainerOrder int
	ContainerRow   ContainerRow
}

type CodeBuilder struct {
	code *store.Code
}

func NewCodeBuilder() CodeBuilder {
	return CodeBuilder{code: &store.Code{Texts: make([]store.DocumentText, 0)}}
}

func (c *CodeBuilder) SetCodeId(id int) *CodeBuilder {
	c.code.Id = id

	return c
}

func (c *CodeBuilder) SetContainerId(id int) *CodeBuilder {
	c.code.Container = id

	return c
}

func (c *CodeBuilder) Push(row CodeRow) {
	if row.Name.Valid && c.code.Name == "" {
		c.code.Name = row.Name.String
	}

	if row.Text.Valid {
		c.code.Texts = append(c.code.Texts, store.DocumentText{
			DocumentId: int(row.SourceId.Int32), Text: row.Text.String, FirstWord: store.WordCoordinate{
				Paragraph: int(row.P1.Int32), Sentence: int(row.S1.Int32), Word: int(row.W1.Int32)},
			LastWord: store.WordCoordinate{
				Paragraph: int(row.P2.Int32),
				Sentence:  int(row.S2.Int32),
				Word:      int(row.W2.Int32),
			},
			Id: int(row.TextId.Int32),
		})
	}
}

func (c *CodeBuilder) Finish() store.Code {
	code := c.code
	c.code = nil

	return *code
}

type ContainerBuilder struct {
	container      *store.CodeContainer
	currentDisplay int
	codeBuilder    *CodeBuilder
}

func NewContainerBuilder(containerId int) ContainerBuilder {
	return ContainerBuilder{
		container: &store.CodeContainer{
			Id:    containerId,
			Codes: make([]store.Code, 0),
		},
		currentDisplay: math.MaxInt64,
	}
}

func (c *ContainerBuilder) Push(row ContainerRow) {
	if row.CodeDisplayOrder != c.currentDisplay {
		c.currentDisplay = row.CodeDisplayOrder

		if c.codeBuilder != nil {
			c.container.Codes = append(c.container.Codes, c.codeBuilder.Finish())
		}

		codeBuilder := NewCodeBuilder()
		c.codeBuilder = codeBuilder.SetContainerId(c.container.Id).SetCodeId(row.CodeId)
	}

	c.codeBuilder.Push(row.CodeRow)
}

func (c *ContainerBuilder) Finish() store.CodeContainer {
	if c.codeBuilder != nil {
		c.container.Codes = append(c.container.Codes, c.codeBuilder.Finish())
	}

	container := c.container

	c.container = nil

	return *container
}

type ContainerListBuilder struct {
	containers       []store.CodeContainer
	containerBuilder *ContainerBuilder
}

func NewContainerListBuilder() ContainerListBuilder {
	return ContainerListBuilder{
		containers: make([]store.CodeContainer, 0, 0),
	}
}

func (c *ContainerListBuilder) Push(row ContainerListRow) {
	if c.containerBuilder != nil && c.containerBuilder.container.Id != row.ContainerId {
		c.containers = append(c.containers, c.containerBuilder.Finish())
		c.containerBuilder = nil
	}

	if c.containerBuilder == nil {
		builder := NewContainerBuilder(row.ContainerId)
		c.containerBuilder = &builder
	}

	c.containerBuilder.Push(row.ContainerRow)
}

func (c *ContainerListBuilder) Finish() []store.CodeContainer {
	if c.containerBuilder != nil {
		c.containers = append(c.containers, c.containerBuilder.Finish())
		c.containerBuilder = nil
	}

	return c.containers
}

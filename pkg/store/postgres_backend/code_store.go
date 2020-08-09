package postgres_backend

import (
	"errors"
	"fmt"
	"github.com/discmonkey/retext/pkg/store"
	"github.com/discmonkey/retext/pkg/version"
	"math"
)

type CodeStore struct {
	db connection
}

func NewCodeStore() (*CodeStore, error) {
	con, err := GetConnection()
	if err != nil {
		return nil, err
	}

	return &CodeStore{db: con}, nil
}

func (c CodeStore) CreateContainer() (store.ContainerId, error) {
	row := c.db.QueryRow(`
		INSERT INTO qode.code_container (display_order) VALUES (0)
		RETURNING id  
	`)

	var id int

	err := row.Scan(&id)

	return id, err
}

func IdDoesNotExistError(objectName string, id int) error {
	return errors.New(fmt.Sprintf("%s with id: <%d> does not exist", objectName, id))
}

func (c CodeStore) CreateCode(name string, containerId store.ContainerId) (store.CodeId, error) {
	row := c.db.QueryRow(`
		
		INSERT INTO qode.code (display_order, name, code_container_id) VALUES (
			(SELECT count(*) FROM qode.code where code_container_id = $1), $2, $1 
		)  
		
		RETURNING id;
	`, containerId, name)

	var id int

	err := row.Scan(&id)

	return id, err
}

func (c CodeStore) CodifyText(codeId store.CodeId, documentId store.FileId, text string, firstWord store.WordCoordinate, lastWord store.WordCoordinate) error {

	// TODO grab parser id from environment variable (or something similar)
	_, err := c.db.Exec(`
		INSERT INTO qode.text (start, stop, value, parser_id, code_id, source_file_id) VALUES 
		(ROW($1, $2, $3), ROW($4, $5, $6), $7, (
		    SELECT id from qode.parser WHERE version = $10
		), $8, $9) 
	`, firstWord.Paragraph, firstWord.Sentence, firstWord.Word,
		lastWord.Paragraph, lastWord.Sentence, lastWord.Word,
		text, codeId, documentId, version.Version)

	return err
}

type codeRow struct {
	Name string

	P1 int
	S1 int
	w1 int

	P2 int
	S2 int
	W2 int

	Text     string
	SourceId int
}

type CodeBuilder struct {
	code store.Code
}

func (c *CodeBuilder) SetId(id int) {
	c.code.Id = id
}

func (c *CodeBuilder) Push(row codeRow) {
	if c.code.Name == "" {
		c.code.Name = row.Name
	}

	if len(row.Text) > 0 {
		c.code.Texts = append(c.code.Texts, store.DocumentText{
			DocumentId: row.SourceId, Text: row.Text, FirstWord: store.WordCoordinate{
				Paragraph: row.P1, Sentence: row.S1, Word: row.w1}, LastWord: store.WordCoordinate{
				Paragraph: row.P2,
				Sentence:  row.S2,
				Word:      row.W2,
			},
		})
	}
}

type ContainerBuilder struct {
	container      store.CodeContainer
	currentDisplay int
	codeBuilder    *CodeBuilder
}

func NewContainerBuilder() ContainerBuilder {
	return ContainerBuilder{
		container:      store.CodeContainer{},
		currentDisplay: math.MaxInt64,
	}
}

func (c *ContainerBuilder) Push(row codeRow, codeDisplayOrder int, codeId int) {
	if codeDisplayOrder != c.currentDisplay {
		if codeDisplayOrder < c.currentDisplay {
			c.container.Main = codeId
		}

		c.currentDisplay = codeDisplayOrder

		c.Finish()

		c.codeBuilder = &CodeBuilder{}
		c.codeBuilder.SetId(codeId)
	}

	c.codeBuilder.Push(row)
}

func (c *ContainerBuilder) Finish() {
	if c.codeBuilder != nil {
		c.container.Codes = append(c.container.Codes, c.codeBuilder.code)
	}
}

func (c CodeStore) GetCode(codeId store.CodeId) (store.Code, error) {
	builder := CodeBuilder{}
	builder.SetId(codeId)

	rows, err := c.db.Query(`
		SELECT code.name, (text.start).paragraph, (text.start).sentence, (text.start).word, 
		       (text.stop).paragraph, (text.stop).sentence, (text.stop).word, text.value, text.source_file_id FROM qode.code code
		LEFT JOIN qode.text text on code.id = text.code_id
		WHERE code.id = $1 
	`, codeId)

	if err != nil {
		return builder.code, err
	}

	row := codeRow{}
	empty := true
	for rows.Next() {
		empty = false
		err = rows.Scan(&row.Name, &row.P1, &row.S1, &row.w1,
			&row.P2, &row.S2, &row.W2, &row.Text, &row.SourceId)

		if err != nil {
			return builder.code, err
		}

		builder.Push(row)
	}

	if empty {
		return builder.code, IdDoesNotExistError("code", codeId)
	}

	return builder.code, nil
}

func (c CodeStore) GetContainer(containerId store.ContainerId) (store.CodeContainer, error) {
	container := store.CodeContainer{
		Main: 0,
	}

	rows, err := c.db.Query(`
		SELECT c.name, c.display_order, c.id, (t.start).paragraph, (t.start).sentence, (t.start).word, 
		       (t.stop).paragraph, (t.stop).sentence, (t.stop).word, t.value, t.source_file_id FROM qode.code c
		LEFT JOIN qode.text t on c.id = t.code_id
		WHERE c.code_container_id = $1
		ORDER BY c.display_order
	`, containerId)

	if err != nil {
		return container, err
	}

	builder := NewContainerBuilder()

	var displayOrder int
	var codeId int
	row := codeRow{}

	for rows.Next() {
		err = rows.Scan(&row.Name, &displayOrder, &codeId, &row.P1, &row.S1, &row.w1,
			&row.P2, &row.S2, &row.W2, &row.Text, &row.SourceId)

		builder.Push(row, displayOrder, codeId)
	}

	builder.Finish()

	builder.container.Id = containerId

	return builder.container, nil
}

func (c CodeStore) GetContainers() ([]store.CodeContainer, error) {
	containers := make([]store.CodeContainer, 0, 0)

	rows, err := c.db.Query(`
		SELECT container.display_order as container_display_order,
		       c.name, c.display_order, c.id, c.code_container_id, (t.start).paragraph, (t.start).sentence, (t.start).word, 
		       (t.stop).paragraph, (t.stop).sentence, (t.stop).word, t.value, t.source_file_id 
		FROM qode.code_container container  
	    LEFT JOIN qode.code c on c.code_container_id = container.id
		LEFT JOIN qode.text t on c.id = t.code_id
		ORDER BY c.code_container_id, c.display_order
	`)

	if err != nil {
		return nil, err
	}

	builder := NewContainerBuilder()
	currentContainer := -1
	containerId := -1
	codeId := -1
	displayOrder := 0
	containerDisplayOrder := 0
	row := codeRow{}

	for rows.Next() {
		err = rows.Scan(&containerDisplayOrder, &row.Name, &displayOrder, &codeId, &containerId, &row.P1, &row.S1, &row.w1,
			&row.P2, &row.S2, &row.W2, &row.Text, &row.SourceId)

		if currentContainer != containerId {
			if currentContainer != -1 {
				containers = append(containers, builder.container)
				builder = NewContainerBuilder()
			}

			currentContainer = containerId
		}

		builder.Push(row, displayOrder, codeId)
		builder.container.Order = containerDisplayOrder
		builder.container.Id = containerId
	}

	if currentContainer != -1 {
		builder.Finish()
		containers = append(containers, builder.container)
	}

	return containers, nil
}

var _ store.CodeStore = CodeStore{}

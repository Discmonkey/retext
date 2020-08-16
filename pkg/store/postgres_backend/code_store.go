package postgres_backend

import (
	"errors"
	"fmt"
	"github.com/discmonkey/retext/pkg/store"
	"github.com/discmonkey/retext/pkg/store/postgres_backend/builders"
	"github.com/discmonkey/retext/pkg/version"
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

func (c CodeStore) CreateContainer(id store.ProjectId) (store.ContainerId, error) {
	row := c.db.QueryRow(`
		INSERT INTO qode.code_container (display_order, project_id) VALUES (0, $1)
		RETURNING id  
	`, id)

	var containerId int

	err := row.Scan(&containerId)

	return containerId, err
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

func (c CodeStore) GetCode(codeId store.CodeId) (store.Code, error) {
	builder := builders.NewCodeBuilder()

	rows, err := c.db.Query(`
		SELECT code.name, code.code_container_id, (text.start).paragraph, (text.start).sentence, (text.start).word, 
		       (text.stop).paragraph, (text.stop).sentence, (text.stop).word, text.value, text.source_file_id FROM qode.code code
		LEFT JOIN qode.text text on code.id = text.code_id
		WHERE code.id = $1 
	`, codeId)

	if err != nil {
		return builder.Finish(), err
	}

	row := builders.CodeRow{}
	var codeContainerId int
	empty := true

	for rows.Next() {
		empty = false
		err = rows.Scan(&row.Name, &codeContainerId, &row.P1, &row.S1, &row.W1,
			&row.P2, &row.S2, &row.W2, &row.Text, &row.SourceId)

		if err != nil {
			return builder.Finish(), err
		}

		builder.Push(row)
	}

	if empty {
		return builder.Finish(), IdDoesNotExistError("code", codeId)
	}

	return builder.SetCodeId(codeId).SetContainerId(codeContainerId).Finish(), nil
}

func (c CodeStore) GetContainer(containerId store.ContainerId) (store.CodeContainer, error) {
	rows, err := c.db.Query(`
		SELECT c.name, c.display_order, c.id, (t.start).paragraph, (t.start).sentence, (t.start).word, 
		       (t.stop).paragraph, (t.stop).sentence, (t.stop).word, t.value, t.source_file_id FROM qode.code c
		LEFT JOIN qode.text t on c.id = t.code_id
		WHERE c.code_container_id = $1
		ORDER BY c.display_order
	`, containerId)

	if err != nil {
		return store.CodeContainer{}, err
	}

	builder := builders.NewContainerBuilder(containerId)

	row := builders.ContainerRow{}

	for rows.Next() {
		err = rows.Scan(&row.CodeRow.Name, &row.CodeDisplayOrder, &row.CodeId, &row.CodeRow.P1,
			&row.CodeRow.S1, &row.CodeRow.W1,
			&row.CodeRow.P2, &row.CodeRow.S2, &row.CodeRow.W2, &row.CodeRow.Text, &row.CodeRow.SourceId)

		builder.Push(row)
	}

	return builder.Finish(), nil
}

func (c CodeStore) GetContainers(id store.ProjectId) ([]store.CodeContainer, error) {

	rows, err := c.db.Query(`
		SELECT container.display_order as container_display_order,
		       c.name, c.display_order, c.id, c.code_container_id, (t.start).paragraph, (t.start).sentence, (t.start).word, 
		       (t.stop).paragraph, (t.stop).sentence, (t.stop).word, t.value, t.source_file_id 
		FROM qode.code_container container  
	    LEFT JOIN qode.code c on c.code_container_id = container.id
		LEFT JOIN qode.text t on c.id = t.code_id
		WHERE container.project_id = $1 
		ORDER BY c.code_container_id, c.display_order
	`, id)

	if err != nil {
		return nil, err
	}

	builder := builders.NewContainerListBuilder()
	row := builders.ContainerListRow{}

	for rows.Next() {
		err = rows.Scan(&row.ContainerOrder, &row.ContainerRow.CodeRow.Name, &row.ContainerRow.CodeDisplayOrder,
			&row.ContainerRow.CodeId, &row.ContainerId, &row.ContainerRow.CodeRow.P1, &row.ContainerRow.CodeRow.S1,
			&row.ContainerRow.CodeRow.W1,
			&row.ContainerRow.CodeRow.P2, &row.ContainerRow.CodeRow.S2, &row.ContainerRow.CodeRow.W2, &row.ContainerRow.CodeRow.Text, &row.ContainerRow.CodeRow.SourceId)

		builder.Push(row)
	}

	return builder.Finish(), nil
}

var _ store.CodeStore = CodeStore{}

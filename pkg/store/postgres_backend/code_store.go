package postgres_backend

import (
	"fmt"
	"github.com/discmonkey/retext/pkg/store"
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

func (c CodeStore) CreateContainer() (store.ContainerID, error) {
	row := c.db.QueryRow(`
		INSERT INTO qode.code_container (display_order) VALUES (0)
		RETURNING id  
	`)

	var id int

	err := row.Scan(&id)

	return id, err
}

func (c CodeStore) CreateCode(name string, containerID store.ContainerID) (store.CodeID, error) {
	row := c.db.QueryRow(`
		
		INSERT INTO qode.code (display_order, name, code_container_id) VALUES (
			(SELECT count(*) FROM qode.code where code_container_id = $1), $2, $1 
		)  
		
		RETURNING id;
	`, containerID, name)

	var id int

	err := row.Scan(&id)

	return id, err
}

func (c CodeStore) CodifyText(codeID store.CodeID, documentID store.FileID, text string, firstWord store.WordCoordinate, lastWord store.WordCoordinate) error {

	// TODO grab parser id from environment variable (or something similar)
	_, err := c.db.Exec(`
		INSERT INTO qode.text (start, stop, value, parser_id, code_id, source_file_id) VALUES 
		(ROW($1, $2, $3), ROW($4, $5, $6), $7, 0, $8, $9) 
	`, firstWord.Paragraph, firstWord.Sentence, firstWord.Word,
		lastWord.Paragraph, lastWord.Sentence, lastWord.Word,
		text, codeID, documentID)

	return err
}

func (c CodeStore) GetCode(codeID store.CodeID) (store.Code, error) {
	code := store.Code{
		ID:    codeID,
		Name:  "",
		Texts: make([]store.DocumentText, 0, 0),
	}

	rows, err := c.db.Query(`
		SELECT c.name, t.start.paragraph, t.start.sentence, t.start.word, 
		       t.stop.paragraph, t.stop.sentence, t.stop.word, t.value, t.source_file_id FROM qode.code c
		LEFT JOIN qode.text t on c.id = t.code_id
		WHERE c.id = $1 
	`, codeID)

	if err != nil {
		return code, err
	}

	var fileID int
	var name string
	var text string

	for rows.Next() {
		coord1 := store.WordCoordinate{}
		coord2 := store.WordCoordinate{}

		err = rows.Scan(&name, coord1.Paragraph, coord1.Sentence, coord1.Word,
			coord2.Paragraph, coord2.Sentence, coord2.Word, &text, &fileID)

		if err != nil {
			return code, err
		}

		d := store.DocumentText{
			DocumentID: fmt.Sprintf("%d", fileID),
			Text:       text,
			FirstWord:  coord1,
			LastWord:   coord2,
		}

	}

	code.Name = name

	return code, nil
}

func (c CodeStore) GetContainer(containerID store.ContainerID) (store.CodeContainer, error) {
	rows, err := c.db.Query(`
		with temp as (SELECT n.id, n.display_order
		into #temp
		FROM qode.code_container n
		WHERE n.id = $1),
	
		SELECT * FROM #temp;
	
		SELECT c.name, t.start, t.stop, t.value, c.display_order
		FROM qode.code c
		INNER JOIN #temp t on c.code_container_id = t.id
		LEFT JOIN qode.text t on c.id = t.code_id

	`, containerID)
}

func (c CodeStore) GetContainers() ([]store.CodeContainer, error) {

}

var _ store.CodeStore = CodeStore{}

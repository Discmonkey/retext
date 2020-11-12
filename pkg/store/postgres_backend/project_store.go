package postgres_backend

import (
	"database/sql"
	"fmt"
	"github.com/discmonkey/retext/pkg/store"
	"time"
)

type ProjectStore struct {
	con *sql.DB
}

func NewProjectStore(con *sql.DB) ProjectStore {
	return ProjectStore{con: con}
}

func (p ProjectStore) CreateProject(name, description string, month, year int) (store.ProjectId, error) {

	t := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	fmt.Println(t)
	query := `
        INSERT INTO qode.project (name, description, created, time_tag) 
        VALUES ($1, $2, now(), $3)
        returning id;
    `

	row := p.con.QueryRow(query, name, description, t)

	var projectId int64

	err := row.Scan(&projectId)

	return projectId, err
}

func (p ProjectStore) GetProject(id store.ProjectId) (store.Project, error) {
	query := `
		SELECT id, name, description, time_tag FROM qode.project WHERE id = $1;
	`

	row := p.con.QueryRow(query, id)

	var project store.Project
	var datetime time.Time
	err := row.Scan(&project.Id, &project.Name, &project.Description, &datetime)

	project.Year = int32(datetime.Year())
	project.Month = int32(datetime.Month())

	return project, err
}

func (p ProjectStore) GetProjects() ([]store.Project, error) {
	query := `
		SELECT id, name, description, time_tag FROM qode.project ORDER BY id 
	`

	results := make([]store.Project, 0)

	rows, err := p.con.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var project store.Project
		var datetime time.Time
		err := rows.Scan(&project.Id, &project.Name, &project.Description, &datetime)
		project.Month = int32(datetime.Month())
		project.Year = int32(datetime.Year())

		if err != nil {
			return nil, err
		}

		results = append(results, project)
	}

	return results, nil
}

var _ store.ProjectStore = &ProjectStore{}

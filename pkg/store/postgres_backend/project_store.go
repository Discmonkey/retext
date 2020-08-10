package postgres_backend

import (
	"database/sql"
	"github.com/discmonkey/retext/pkg/store"
)

type ProjectStore struct {
	con *sql.DB
}

func NewProjectStore() (ProjectStore, error) {
	con, err := GetConnection()
	if err != nil {
		return ProjectStore{}, err
	}

	return ProjectStore{con: con}, nil
}

func (p ProjectStore) CreateProject(name, description string, month, year int) (store.ProjectId, error) {
	query := `
        INSERT INTO qode.project (name, description, created, time_tag) 
        VALUES ($1, $2, now(), format('%s-%s-%s', $3, $4, 0)::time)
    `

	row := p.con.QueryRow(query, name, description, month, year)

	var projectId int

	err := row.Scan(&projectId)

	return projectId, err
}

func (p ProjectStore) GetProject(id store.ProjectId) (store.Project, error) {
	panic("implement me")
}

func (p ProjectStore) GetProjects() ([]store.Project, error) {
	panic("implement me")
}

var _ store.ProjectStore = &ProjectStore{}

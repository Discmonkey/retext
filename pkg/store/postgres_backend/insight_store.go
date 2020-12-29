package postgres_backend

import (
	"database/sql"
	"github.com/discmonkey/retext/pkg/store"
	"github.com/discmonkey/retext/pkg/store/postgres_backend/builders"
	"github.com/lib/pq"
)

type InsightStore struct {
	db connection
}

func NewInsightStore(db *sql.DB) InsightStore {
	return InsightStore{db: db}
}

func (s InsightStore) CreateInsight(projectId store.ProjectId, insightText string, textIds []store.TextId) (store.InsightId, error) {
	row := s.db.QueryRow(`
		WITH newInsight as (
			INSERT INTO qode.insight (project_id, value) VALUES ($1, $2) RETURNING id
		)
		INSERT INTO qode.insight_text (insight_id, text_id)
		SELECT id, unnest($3::integer[]) FROM newInsight RETURNING insight_id
	`, projectId, insightText, pq.Array(textIds))

	var insightId int64

	err := row.Scan(&insightId)

	if err != nil {
		return 0, err
	}

	return insightId, err
}

func (s InsightStore) GetInsights(projectId store.ProjectId) ([]store.Insight, error) {
	rows, err := s.db.Query(`
		SELECT i.id, i.value, it.text_id FROM qode.insight i
		INNER JOIN qode.insight_text it on i.id = it.insight_id
		WHERE i.project_id = $1
		ORDER BY i.id
	`, projectId)

	if err != nil {
		return nil, err
	}

	builder := builders.NewInsightListBuilder()
	row := builders.InsightListRow{}

	for rows.Next() {
		err = rows.Scan(&row.InsightId, &row.Value, &row.TextId)

		builder.Push(row)
	}

	return builder.Finish(), err
}

func (s InsightStore) DeleteInsightText(insightId store.InsightId, textId store.TextId) error {
	_, err := s.db.Query("DELETE FROM qode.insight_text WHERE insight_id = $1 AND text_id = $2", insightId, textId)

	if err != nil {
		return err
	}

	return nil
}

var _ store.InsightStore = InsightStore{}

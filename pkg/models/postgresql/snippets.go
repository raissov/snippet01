package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"raissov/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *pgxpool.Conn
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES($1, $2, now(), now() + make_interval(days =>$3)) RETURNING id`

	var id int
	err := m.DB.QueryRow(context.Background(), stmt, title, content, expires).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil

}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > now() AND id = $1`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	sm := &models.Snippet{}

	err := row.Scan(&sm.ID, &sm.Title, &sm.Content, &sm.Created, sm.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}

	}
	return sm, nil

}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires, FROM snippets WHERE expires > now() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(context.Background(), stmt)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, s.Content, s.Created, s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

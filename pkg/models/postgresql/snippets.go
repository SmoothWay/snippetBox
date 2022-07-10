package postgresql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/SmoothWay/snippetBox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	var id int
	stmt := fmt.Sprintf(`INSERT INTO snippets (title, content, created, expires)
		VALUES($1, $2, now(), now() + '%s days') RETURNING id`, expires)

	err := m.DB.QueryRow(stmt, title, content).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	s := &models.Snippet{}
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE id=$1`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	s.Created, _ = time.Parse("2006-01-02 15:04", s.Created.String()[:16])
	s.Expires, _ = time.Parse("2006-01-02 15:04", s.Expires.String()[:16])
	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
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

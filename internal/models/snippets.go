package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`
	item := m.DB.QueryRow(query, id)
	s := &Snippet{}
	err := item.Scan(&s.ID, &s.Title, &s.Content, &s.Expires, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	query := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() LIMIT 10"
	items, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer m.DB.Close()
	snippets := []*Snippet{}
	for items.Next() {
		s := &Snippet{}
		err := items.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}
		snippets = append(snippets, s)
	}
	return snippets, nil
}

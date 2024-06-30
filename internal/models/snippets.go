package models

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// Snippet Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Snippet struct {
	ID      uint64
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DBPool *pgxpool.Pool
}

// Insert This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) error {
	expTime := time.Now().Add(time.Duration(expires) * time.Hour * 24)
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES ($1, $2, CURRENT_TIMESTAMP, $3)`
	ct, err := m.DBPool.Exec(context.Background(), stmt, title, content, expTime)
	if err != nil {
		return err
	} else if ct.RowsAffected() <= 0 {
		return errors.New("no insertions were made, unexpected behaviour")
	}
	return nil
}

// Get This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	s := &Snippet{}
	stmt := `SELECT id, title, content, created, expires 
			FROM snippets
			WHERE expires > CURRENT_TIMESTAMP AND id=$1`
	r := m.DBPool.QueryRow(context.Background(), stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	err := r
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNoRows
	}
	return s, nil
}

// Latest This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT * FROM snippets WHERE expires > CURRENT_TIMESTAMP ORDER BY id DESC LIMIT 10`
	rows, _ := m.DBPool.Query(context.Background(), stmt)
	defer rows.Close()
	var snippets []*Snippet
	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err := rows.Err(); err != nil { // Preferred error handling on pgx.rows
		return nil, err
	}
	return snippets, nil
}

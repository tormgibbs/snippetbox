package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)



type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *pgxpool.Pool
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 DAY' * $3)
	RETURNING id`
	
	// Use the QueryRow method to execute the SQL statement.
	// Pass in the args for the placeholders, and scan the returned id into a
	// variable.
	var id int
	err := m.DB.QueryRow(context.Background(), stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {

	// Initialize a pointer to a new Snippet struct.
	s := &Snippet{}

	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > CURRENT_TIMESTAMP AND id = $1`

	// Use the QueryRow method to execute the SQL statement.
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	// Scan the values from the row into the Snippet struct.
	// err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	// Return the Snippet struct.
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > CURRENT_TIMESTAMP
	ORDER BY created DESC
	LIMIT 10`

	// Use the Query method to execute the SQL statement.
	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}

	// We need to close the rows when we're done with them.
	defer rows.Close()

	// Initialize an empty slice to store the Snippet structs.
	snippets := []*Snippet{}

	// Iterate over the rows in the resultset.
	for rows.Next() {
		// Initialize a pointer to a new Snippet struct.
		s := &Snippet{}

		// Scan the values from the row into the Snippet struct.
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// Append the Snippet struct to the slice.
		snippets = append(snippets, s)
	}

	// When the rows are closed, check for any errors.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return the slice of Snippet structs.
	return snippets, nil
}


// func (m *SnippetModel) Insert(title string, content string, expires int) (int,
// 	error) {
// 	stmt := `INSERT INTO snippets (title, content, created, expires)
// 	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
// 	result, err := m.DB.Exec(stmt, title, content, expires)
// 	if err != nil {
// 	return 0, err
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 	return 0, err
// 	}
// 	return int(id), nil
// }
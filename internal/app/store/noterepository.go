package store

import (
	"github.com/AlexCorn999/notes/internal/app/note"
)

type NoteRepository struct {
	store *Store
}

// CreateNote creates a note
func (n *NoteRepository) CreateNote(note *note.Note, user_id int) error {
	if _, err := n.store.db.Exec(
		"insert into notes (title, content, user_id) values ($1, $2, $3)",
		note.Title,
		note.Text,
		user_id,
	); err != nil {
		return err
	}
	return nil
}

// GetList displays all user notes
func (n *NoteRepository) GetList(user_id int) ([]note.Note, error) {
	rows, err := n.store.db.Query("select title from notes where user_id = $1", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	titles := make([]note.Note, 0)
	for rows.Next() {
		var n note.Note
		if err := rows.Scan(&n.Title); err != nil {
			return nil, err
		}
		titles = append(titles, n)
	}

	if rows.Err(); err != nil {
		return nil, err
	}

	return titles, nil
}

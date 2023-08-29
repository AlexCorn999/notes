package note

import "errors"

type Note struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// Validate validates notes
func (n *Note) Validate() error {
	if len(n.Title) < 1 {
		return errors.New("no title")
	}
	if len(n.Text) < 1 {
		return errors.New("no text")
	}
	return nil
}

package yandex

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/AlexCorn999/notes/internal/app/note"
)

type SpellingError struct {
	Code    int      `json:"code"`
	Pos     int      `json:"pos"`
	Row     int      `json:"row"`
	Col     int      `json:"col"`
	Len     int      `json:"len"`
	Word    string   `json:"word"`
	Suggest []string `json:"s"`
}

type SpellingResponse []SpellingError

// CheckAll scans the note for errors
func CheckAll(n *note.Note) error {

	spellingErrorsTitle, err := checkSpelling(n.Title)
	if err != nil {
		return err
	}

	for _, err := range spellingErrorsTitle {
		if len(err.Suggest) > 0 {
			n.Title = replaceWordInText(n.Title, err.Word, err.Suggest[0])
		}
	}

	spellingErrorsText, err := checkSpelling(n.Text)
	if err != nil {
		return err
	}

	for _, err := range spellingErrorsText {
		if len(err.Suggest) > 0 {
			n.Text = replaceWordInText(n.Text, err.Word, err.Suggest[0])
		}
	}

	return nil
}

// checkSpelling checks for errors
func checkSpelling(text string) ([]SpellingError, error) {
	urll := "https://speller.yandex.net/services/spellservice.json/checkText"
	params := url.Values{}
	params.Set("text", text)
	params.Set("lang", "ru,en")
	params.Set("options", "0")
	params.Set("format", "plain")

	resp, err := http.Get(urll + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var spellingResponse SpellingResponse
	err = json.Unmarshal(body, &spellingResponse)
	if err != nil {
		return nil, err
	}

	return spellingResponse, nil
}

// replaceWordInText replaces the words
func replaceWordInText(text, oldWord, newWord string) string {
	replacedText := strings.ReplaceAll(text, oldWord, newWord)
	return replacedText
}

package postgres

import (
	"errors"

	"github.com/chuabingquan/snippets"
	"github.com/jmoiron/sqlx"
)

// SnippetService implements the snippets.SnippetService interface
type SnippetService struct {
	DB *sqlx.DB
}

// Snippet queries the database and returns a snippets.Snippet instance with the
// given snippetID should it exist and belong to the user with the given userID
func (ss SnippetService) Snippet(userID string, snippetID string) (snippets.Snippet, error) {
	var snippet snippets.Snippet
	err := ss.DB.QueryRowx("SELECT * FROM snippet WHERE id=$1 AND account_id=$2", snippetID, userID).StructScan(&snippet)
	if err != nil {
		return snippet, errors.New("Error retrieving snippet: " + err.Error())
	}
	return snippet, nil
}

// Snippets queries the database and returns a slice of snippets.Snippet given
// a userID they associate with
func (ss SnippetService) Snippets(userID string) ([]snippets.Snippet, error) {
	var snippetSlice []snippets.Snippet
	rows, err := ss.DB.Queryx("SELECT * FROM snippet WHERE account_id=$1", userID)
	if err != nil {
		return nil, errors.New("Error retrieving snippets: " + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var snippet snippets.Snippet
		err := rows.StructScan(&snippet)
		if err != nil {
			return nil, errors.New("Error retrieving snippets: " + err.Error())
		}
		snippetSlice = append(snippetSlice, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New("Error retrieving snippets: " + err.Error())
	}

	return snippetSlice, nil
}

// CreateSnippet inserts a new snippet into the database for a given userID
func (ss SnippetService) CreateSnippet(s snippets.Snippet) error {
	_, err := ss.DB.NamedExec(`INSERT INTO snippet(account_id, filename, description, is_public) VALUES(:account_id, :filename, :description, :is_public)`, s)
	if err != nil {
		return errors.New("Error creating snippet: " + err.Error())
	}
	return nil
}

// UpdateSnippet updates an existing snippet in the database
func (ss SnippetService) UpdateSnippet(updatedSnippet snippets.Snippet) error {
	res, err := ss.DB.NamedExec(`UPDATE snippet SET account_id=:account_id, filename=:filename, description=:description,
								is_public=:is_public WHERE id=:id`, updatedSnippet)
	if err != nil {
		return errors.New("Error updating snippet: " + err.Error())
	}
	if rows, err := res.RowsAffected(); err != nil {
		return errors.New("Error checking rows affected after snippet update: " + err.Error())
	} else if rows < 1 {
		return errors.New("Snippet with the given UUID does not exist")
	}
	return nil
}

// DeleteSnippet removes a snippet from the database should its given snippetID
// exist and is associated with the user with the given userID
func (ss SnippetService) DeleteSnippet(userID string, snippetID string) error {
	_, err := ss.DB.Exec("DELETE FROM snippet WHERE id=$1 AND account_id=$2", snippetID, userID)
	if err != nil {
		return errors.New("Error deleting snippet: " + err.Error())
	}
	return nil
}

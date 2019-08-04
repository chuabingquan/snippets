package postgres

import (
	"database/sql"
	"errors"

	"github.com/chuabingquan/snippets"
	"github.com/jmoiron/sqlx"
)

// AuthenticationService implements the snippets.AuthenticationService interface
type AuthenticationService struct {
	DB            *sqlx.DB
	HashUtilities snippets.HashUtilities
}

// Authenticate queries the database and verifies a user's credentials
func (as AuthenticationService) Authenticate(username string, password string) (bool, error) {
	var passwordHash string
	err := as.DB.QueryRowx("SELECT password_hash FROM account WHERE username=$1", username).Scan(&passwordHash)
	if err == sql.ErrNoRows {
		return false, nil // no such username exists
	} else if err != nil {
		// any other error represents a failure
		return false, errors.New("Database query failed for authentication: " + err.Error())
	}
	return as.HashUtilities.CompareHashWithString(passwordHash, password), nil
}

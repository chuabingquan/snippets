package postgres

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

// Open returns a connection to the database given relevant credentials
func Open(dbURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, errors.New("Could not connect to database: " + err.Error())
	}
	return db, nil
}

// DBUrl represents the structure of a database connection string
type DBUrl struct {
	Protocol string
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	Sslmode  string
}

// GetURL returns a database connection string constructed from values supplied by the DBUrl struct
func (du DBUrl) GetURL() string {
	return du.Protocol + "://" + du.User + ":" + du.Password + "@" + du.Host + ":" + du.Port + "/" + du.Name + "?sslmode=" + du.Sslmode
}

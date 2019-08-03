package postgres

import (
	"errors"

	"github.com/chuabingquan/snippets"
	"github.com/jmoiron/sqlx"
)

// UserService implements the snippets.UserService interface
type UserService struct {
	DB            *sqlx.DB
	HashUtilities snippets.HashUtilities
}

// User returns a snippets.User after querying from the database given a userID,
// else, an error occurs such as when the user isn't found
func (us UserService) User(userID string) (snippets.User, error) {
	var user snippets.User
	err := us.DB.QueryRowx("SELECT * FROM account WHERE id=$1", userID).StructScan(&user)
	if err != nil {
		return user, errors.New("Error retrieving user: " + err.Error())
	}
	return user, nil
}

// UserByUsername performs the same operation as User but takes in a username instead
// of a userID as an argument
func (us UserService) UserByUsername(username string) (snippets.User, error) {
	var user snippets.User
	err := us.DB.QueryRowx("SELECT * FROM account WHERE username=$1", username).StructScan(&user)
	if err != nil {
		return user, errors.New("Error retrieving user: " + err.Error())
	}
	return user, nil
}

// Users returns all users from the database in the form of a snippets.User slice
func (us UserService) Users() ([]snippets.User, error) {
	users := []snippets.User{}
	rows, err := us.DB.Queryx("SELECT * FROM account")
	if err != nil {
		return nil, errors.New("Error retrieving users: " + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var user snippets.User
		err := rows.StructScan(&user)
		if err != nil {
			return nil, errors.New("Error retrieving users: " + err.Error())
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New("Error retrieving users: " + err.Error())
	}

	return users, nil
}

// CreateUser inserts the data from a given snippets.User instance into the database
func (us UserService) CreateUser(u snippets.User) error {
	hash, err := us.HashUtilities.HashAndSalt(u.Password)
	if err != nil {
		return errors.New("Error creating user: " + err.Error())
	}
	u.PasswordHash = hash

	_, err = us.DB.NamedExec(`INSERT INTO account(email, username, password_hash, first_name, last_name) VALUES(:email, :username, :password_hash, :first_name, :last_name)`, u)
	if err != nil {
		return errors.New("Error creating user: " + err.Error())
	}
	return nil
}

// UpdateUser takes in a snippets.User instance and updates the relevant database user record accordingly
func (us UserService) UpdateUser(updatedUser snippets.User) error {
	// only create new password hash if user updates their password (password pointer field not nil)
	if updatedUser.Password != "" {
		hash, err := us.HashUtilities.HashAndSalt(updatedUser.Password)
		if err != nil {
			return errors.New("Error updating user: " + err.Error())
		}
		updatedUser.PasswordHash = hash
	}

	res, err := us.DB.NamedExec(`UPDATE account SET email=:email, username=:username, password_hash=:password_hash, 
					first_name=:first_name, last_name=:last_name WHERE id=:id`, updatedUser)
	if err != nil {
		return errors.New("Error updating user: " + err.Error())
	}
	if rows, err := res.RowsAffected(); err != nil {
		return errors.New("Error checking rows affected after user update: " + err.Error())
	} else if rows < 1 {
		return errors.New("User with the given UUID does not exist")
	}
	return nil
}

// DeleteUser removes a user with a matching userID (given) from the database
func (us UserService) DeleteUser(userID string) error {
	_, err := us.DB.Exec("DELETE FROM account WHERE id=$1", userID)
	if err != nil {
		return errors.New("Error deleting user: " + err.Error())
	}
	return nil
}

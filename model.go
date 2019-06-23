package snippets

// User represents a registered person of this application who can create snippets
type User struct {
	ID           string `json:"userId" db:"id"`
	Email        string `json:"email" db:"email"`
	Username     string `json:"username" db:"username"`
	PasswordHash string `json:"-" db:"password_hash"`
	PasswordSalt string `json:"-" db:"password_salt"`
	FirstName    string `json:"firstName" db:"first_name"`
	LastName     string `json:"lastName" db:"last_name"`
	// Created/Updated datetime
}

// UserService provides a set of operations that can be applied on the User struct
type UserService interface {
	User(userID string) (User, error)
	Users() ([]User, error)
	CreateUser(u User) error
	UpdateUser(userID string, updatedUser User) error
	DeleteUser(userID string) error
}

// Snippet represents a piece of code published by a user
type Snippet struct {
	ID          string `json:"snippetId" db:"id"`
	Filename    string `json:"filename" db:"filename"`
	Description string `json:"description" db:"description"`
	Public      bool   `json:"isPublic" db:"is_public"`
	// Created/Updated datetime
}

// SnippetService provides a set of operations that can be applied to the Snippet struct
type SnippetService interface {
	Snippet(UserID string, snippetID string) (Snippet, error)
	Snippets(UserID string) ([]Snippet, error)
	CreateSnippet(UserID string, s Snippet) error
	UpdateSnippet(UserID string, snippetID string, updatedSnippet Snippet) error
	DeleteSnippet(UserID string, snippetID string) error
}

package snippets

// User represents a registered person of this application who can create snippets
type User struct {
	ID           string `json:"userId" db:"id"`
	Email        string `json:"email" db:"email"`
	Username     string `json:"username" db:"username"`
	Password     string `json:"password,omitempty"`
	PasswordHash string `json:"-" db:"password_hash"`
	FirstName    string `json:"firstName" db:"first_name"`
	LastName     string `json:"lastName" db:"last_name"`
	// Created/Updated datetime
}

// UserService provides a set of operations that can be applied on the User struct
type UserService interface {
	User(userID string) (User, error)
	Users() ([]User, error)
	CreateUser(u User) error
	UpdateUser(updatedUser User) error
	DeleteUser(userID string) error
}

// Snippet represents a piece of code published by a user
type Snippet struct {
	ID          string `json:"snippetId" db:"id"`
	Filename    string `json:"filename" db:"filename"`
	Description string `json:"description" db:"description"`
	Public      bool   `json:"isPublic" db:"is_public"`
	Owner       string `json:"-" db:"account_id"`
	// Created/Updated datetime
}

// SnippetService provides a set of operations that can be applied to the Snippet struct
type SnippetService interface {
	Snippet(userID string, snippetID string) (Snippet, error)
	Snippets(userID string) ([]Snippet, error)
	CreateSnippet(s Snippet) error
	UpdateSnippet(updatedSnippet Snippet) error
	DeleteSnippet(userID string, snippetID string) error
}

// HashUtilities provides a set of operations relating to hashing and hash comparisons
type HashUtilities interface {
	HashAndSalt(s string) (string, error)
	CompareHashWithString(hash string, s string) bool
}

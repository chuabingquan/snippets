package snippets

// User represents a registered person of this application who can create snippets
type User struct {
	ID           string `json:"userId"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
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
	ID          string `json:"snippetId"`
	Filename    string `json:"filename"`
	Description string `json:"description"`
	Public      bool   `json:"isPublic"`
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

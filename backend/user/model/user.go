package model

const (
	FullNameMaxlength = 32
	FullNameMinLength = 8
	PasswordMaxLength = 16
	PasswordMinLength = 8
)

// AuthInput contains fields for user authentication.
type AuthInput struct {
	PublicID `json:"id"`
	Password string `json:"password"`
}

// AuthOutput contains the auth token returned after successful authentication.
type AuthOutput struct {
	Token string `json:"token"`
}

// CreateInput contains fields for user registration.
type CreateInput struct {
	FullName `json:"full_name"`
	Password string `json:"password"`
}

// CreateOutput contains fields returned after a successful user registration.
type CreateOutput struct {
	ID
	PublicID `json:"id"`
}

// ID is a database primary key type.
type ID uint

// PublicID is the externally visible user ID.
type PublicID string

// FullName represents the user name.
type FullName string

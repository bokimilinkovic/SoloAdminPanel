package dto

// User is a view representation of domain.User.
type User struct {
	ID        uint   `json:"id"`
	Picture   string `json:"picture"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

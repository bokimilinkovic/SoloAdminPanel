package google

import (
	"github.com/kolosek/pkg/model/domain"
)

type Claims struct {
	Subject    string `json:"sub"`     // An identifier for the user, unique among all Google accounts.
	Name       string `json:"name"`    // The user's full name, in a displayable form.
	PictureURL string `json:"picture"` // The URL of the user's profile picture.
	Email      string `json:"email"`   // The preferred email of the user.
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

// ToUser transforms information obtained from an ID token claims to domain.User.
func (c *Claims) ToUser() *domain.User {
	return &domain.User{
		FirstName: c.GivenName,
		ImageURL:  c.PictureURL,
		Email:     c.Email,
		LastName:  c.FamilyName,
	}
}

package cookie

import (
	"net/http"
	"time"
)

// Generator holds methods for creating http cookies.
type Generator struct {
	domain   string
	lifetime time.Duration
}

// NewGenerator creates new cookie generator.
func NewGenerator(domain string, lifetime time.Duration) *Generator {
	return &Generator{
		domain:   domain,
		lifetime: lifetime,
	}
}

// HTTPOnly creates a new httponly cookie.
func (g *Generator) HTTPOnly(name, value, path string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Path:     path,
		Domain:   g.domain,
		Expires:  time.Now().Add(g.lifetime).UTC(),
	}
}

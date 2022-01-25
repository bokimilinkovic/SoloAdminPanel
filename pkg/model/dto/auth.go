package dto

// AuthRequest represents model for sending google oauth request.
type AuthRequest struct {
	AuthCode string `json:"auth_code"`
}

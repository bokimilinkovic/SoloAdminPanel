package google

// VerificationError represents an error that occurred during OAuth token verification.
type VerificationError struct {
	err error
}

// Error returns the underlying error message.
func (e *VerificationError) Error() string {
	msg := e.err.Error()
	return msg
}

// ExchangeError represents an error that occurred during authorization code exchange with Google.
type ExchangeError struct {
	err error
}

// Error returns the underlying error message.
func (e *ExchangeError) Error() string {
	msg := e.err.Error()
	return msg
}

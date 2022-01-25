package domain

import "time"

// User is a domain representation of a user.
type User struct {
	ID             int64 `gorm:"primaryKey"`
	GUID           string
	Email          string
	FirstName      string
	MiddleName     string
	LastName       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ImageURL       string
	Score          int
	BorrowingLimit int64
	Birthday       *time.Time
	BlockStatus    BlockStatus
	BlockReason    *string
}

// BlockStatus defines possible custom assigned status for users.
type BlockStatus int

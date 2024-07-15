package domain

import (
	"time"

	"github.com/google/uuid"
)

// User is an entity that represents a user
type User struct {
	ID        uuid.UUID
	Login     string
	Password  string
	CreatedAt time.Time
}

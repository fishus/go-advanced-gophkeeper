package domain

import (
	"time"

	"github.com/google/uuid"
)

// User Пользователь
type User struct {
	ID        uuid.UUID
	Login     string
	Password  string
	CreatedAt time.Time
}

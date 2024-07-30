package domain

import (
	"github.com/google/uuid"
)

// TokenPayload is an entity that represents the payload of the token
type TokenPayload struct {
	UserID uuid.UUID
}

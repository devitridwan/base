package entity

import (
	"time"

	"github.com/google/uuid"
)

type GetUserResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	StartDate time.Time `json:"username"`
}

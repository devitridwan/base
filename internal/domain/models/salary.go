package models

import (
	"github.com/google/uuid"
)

type User struct {
	UserID   uuid.UUID `db:"user_id"`
	Username uuid.UUID `db:"username"`
	Password uuid.UUID `db:"password"`
}

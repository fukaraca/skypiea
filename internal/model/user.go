package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  int
	UUID                uuid.UUID
	Firstname, Lastname string
	Email               string
	Role, Status        string
	CreatedAt           time.Time
}

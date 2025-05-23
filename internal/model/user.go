package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleAdmin   = "admin"
	RoleUserStd = "user_std"
	RoleUserVip = "user_vip"

	StatusNew = "New"
)

type User struct {
	ID                  int
	UUID                uuid.UUID
	Firstname, Lastname string
	Email               string
	Role, Status        string
	CreatedAt           time.Time
}

type SessionCookie struct {
	Name, Value, Path, Domain string
	MaxAge                    int
	Secure, HTTPOnly          bool
}

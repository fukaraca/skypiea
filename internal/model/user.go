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
	Email, PhoneNumber  string
	Role, Status        string
	AboutMe, Summary    string
	CreatedAt           time.Time
}

type SessionCookie struct {
	Name, Value, Path, Domain string
	MaxAge                    int
	Secure, HTTPOnly          bool
}

package entity

import (
	"time"

	"cloud.google.com/go/datastore"
)

type User struct {
	UID          string
	Gpa          float64
	WishLab      *string
	ConfirmedLab *string
	Year         int
	Role         Role
	Reason       string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleAudience Role = "audience"
)

var RoleByValue = map[string]Role{
	"admin":    RoleAdmin,
	"audience": RoleAudience,
}

const KindUser = "user"

func NewUserKey(uid string) *datastore.Key {
	return datastore.NameKey(KindUser, uid, nil)
}

func NewUser(uid string, gpa float64, year int, role Role, createdAt time.Time) (*User, *datastore.Key) {
	return &User{
		UID:       uid,
		Gpa:       gpa,
		Year:      year,
		Role:      role,
		CreatedAt: createdAt,
	}, NewUserKey(uid)
}

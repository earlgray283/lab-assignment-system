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
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

const KindUser = "user"

func NewUserKey(uid string) *datastore.Key {
	return datastore.NameKey(KindUser, uid, nil)
}

func NewUser(uid string, gpa float64, year int, createdAt time.Time) (*User, *datastore.Key) {
	return &User{
		UID:       uid,
		Gpa:       gpa,
		Year:      year,
		CreatedAt: createdAt,
	}, NewUserKey(uid)
}

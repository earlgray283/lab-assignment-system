package repository

import (
	"time"

	"cloud.google.com/go/datastore"
)

type User struct {
	UID           string
	Email         string
	StudentNumber string
	Name          string
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

const KindUser = "user"

func NewUserKey(uid string) *datastore.Key {
	return datastore.NameKey(KindUser, uid, nil)
}

func NewUser(uid, email, studentNumber, name string) (*User, *datastore.Key) {
	return &User{
		UID:           uid,
		Email:         email,
		StudentNumber: studentNumber,
		Name:          name,
		CreatedAt:     time.Now(),
	}, NewUserKey(uid)
}

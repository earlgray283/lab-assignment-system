package repository

import (
	"time"

	"cloud.google.com/go/datastore"
)

type User struct {
	UID          string
	Gpa          float64
	Lab1         *string
	Lab2         *string
	Lab3         *string
	ConfirmedLab string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

const KindUser = "user"

func NewUserKey(uid string) *datastore.Key {
	return datastore.NameKey(KindUser, uid, nil)
}

func NewUser(uid, lab1, lab2, lab3 string, gpa float64, createdAt time.Time) (*User, *datastore.Key) {
	return &User{
		UID:       uid,
		Lab1:      &lab1,
		Lab2:      &lab2,
		Lab3:      &lab3,
		Gpa:       gpa,
		CreatedAt: createdAt,
	}, NewUserKey(uid)
}

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
	Lab1          string
	Lab2          string
	Lab3          string
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

const KindUser = "user"

func NewUserKey(uid string) *datastore.Key {
	return datastore.NameKey(KindUser, uid, nil)
}

func NewUser(uid, email, studentNumber, name, lab1, lab2, lab3 string) (*User, *datastore.Key) {
	return &User{
		UID:           uid,
		Email:         email,
		StudentNumber: studentNumber,
		Name:          name,
		Lab1:          lab1,
		Lab2:          lab2,
		Lab3:          lab3,
		CreatedAt:     time.Now(),
	}, NewUserKey(uid)
}

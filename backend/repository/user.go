package repository

import (
	"time"

	"cloud.google.com/go/datastore"
)

type User struct {
	UID           string
	Email         string
	StudentNumber int
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

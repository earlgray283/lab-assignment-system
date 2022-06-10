package repository

import (
	"time"

	"cloud.google.com/go/datastore"
)

type GradeRequestToken struct {
	UID       string    `json:"uid,omitempty"`
	Token     string    `json:"token,omitempty"`
	Expires   time.Time `json:"expires,omitempty"`
	CreatedAt time.Time `json:"-"`
}

const KindGradeRequestToken = "gradeRequestToken"

func NewGradeRequestTokenKey(token string) *datastore.Key {
	return datastore.NameKey(KindGradeRequestToken, token, nil)
}

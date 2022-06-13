package repository

import (
	"time"

	"cloud.google.com/go/datastore"
)

type Grade struct {
	UID           string    `json:"uid,omitempty"`
	StudentName   string    `json:"studentName,omitempty"`
	StudentNumber string    `json:"studentNumber,omitempty"`
	Gpa           float64   `json:"gpa,omitempty"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
}

const KindGrade = "grade"

func NewGradeKey(uid string) *datastore.Key {
	return datastore.NameKey(KindGrade, uid, nil)
}

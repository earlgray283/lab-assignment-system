package repository

import (
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
)

type Grade struct {
	UnitSum       int          `json:"unitSum,omitempty"`
	GpSum         int          `json:"gpSum,omitempty"`
	Gpa           float64      `json:"gpa,omitempty"`
	Grades        []ChildGrade `json:"grades,omitempty"`
	StudentName   string       `json:"studentName,omitempty"`
	StudentNumber int          `json:"studentNumber,omitempty"`

	Year      int       `json:"year,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

type ChildGrade struct {
	UnitNum    int     `json:"unitNum,omitempty"`
	Gp         float64 `json:"gp,omitempty"`
	ReportedAt string  `json:"reportedAt,omitempty"`
}

const KindGrade = "grade"

func NewGradesKey(studentNumber int) *datastore.Key {
	return datastore.NameKey(KindGrade, fmt.Sprintf("%d", studentNumber), nil)
}

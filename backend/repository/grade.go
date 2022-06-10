package repository

import (
	"strconv"

	"cloud.google.com/go/datastore"
)

type Grade struct {
	SubjectGrades []SubjectGrade `json:"subjectGrades,omitempty"`
	StudentName   string         `json:"studentName,omitempty"`
	StudentNumber int            `json:"studentNumber,omitempty"`
}

type SubjectGrade struct {
	SubjectName string  `json:"subjectName,omitempty"` // 科目名
	UnitNum     int     `json:"unitNum,omitempty"`     // 単位
	Point       int     `json:"point,omitempty"`       // 点数
	Gp          float64 `json:"gp,omitempty"`          // GP
	ReportedAt  string  `json:"reportedAt,omitempty"`  // 報告日
}

const KindGrade = "grade"

func NewGradeKey(studentNumber int) *datastore.Key {
	return datastore.NameKey(KindGrade, strconv.Itoa(studentNumber), nil)
}

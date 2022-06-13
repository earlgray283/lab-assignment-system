package models

type SubjectGrade struct {
	SubjectName string  `json:"subjectName,omitempty"` // 科目名
	UnitNum     int     `json:"unitNum,omitempty"`     // 単位
	Point       int     `json:"point,omitempty"`       // 点数
	Gp          float64 `json:"gp,omitempty"`          // GP
	ReportedAt  string  `json:"reportedAt,omitempty"`  // 報告日
}

type Grade struct {
	Grades        []SubjectGrade `json:"grades,omitempty"`
	StudentName   string
	StudentNumber string
}

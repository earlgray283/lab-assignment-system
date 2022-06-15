package models

type SubjectGrade struct {
	UnitNum    int     `json:"unitNum,omitempty"`    // 単位
	Point      float64 `json:"point,omitempty"`      // 点数
	Gp         float64 `json:"gp,omitempty"`         // GP
	ReportedAt string  `json:"reportedAt,omitempty"` // 報告日
}

type Grade struct {
	Grades        []SubjectGrade `json:"grades,omitempty"`
	StudentName   string         `json:"studentName,omitempty"`
	StudentNumber string         `json:"studentNumber,omitempty"`
}

package lib

import (
	"lab-assignment-system-backend/repository"
	"time"
)

type CalculateGpaOption struct {
	Until             time.Time
	ExcludeLowerPoint int
}

func CalculateGpa(grade *repository.Grade, option *CalculateGpaOption) float64 {
	gpSum, unitNum := 0.0, 0
	for _, subjectGrade := range grade.SubjectGrades {
		if subjectGrade.Point < option.ExcludeLowerPoint {
			continue
		}
		// TODO:
		// if subjectGrade.ReportedAt < option.Until {
		// 	continue
		// }
		gpSum += subjectGrade.Gp * float64(subjectGrade.UnitNum)
		unitNum += subjectGrade.UnitNum
	}
	return gpSum / float64(unitNum)
}

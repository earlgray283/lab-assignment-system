package lib

import (
	"lab-assignment-system-backend/server/models"
	"time"
)

type CalculateGpaOption struct {
	Until             time.Time
	ExcludeLowerPoint int
}

func CalculateGpa(subjectGrades []models.SubjectGrade, option *CalculateGpaOption) float64 {
	gpSum, unitNum := 0.0, 0
	for _, subjectGrade := range subjectGrades {
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

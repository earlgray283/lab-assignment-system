package lib

import (
	"lab-assignment-system-backend/server/models"
	"log"
	"time"
)

type CalculateGpaOption struct {
	Until             time.Time
	ExcludeLowerPoint float64
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
	log.Println(gpSum, unitNum)
	return gpSum / float64(unitNum)
}

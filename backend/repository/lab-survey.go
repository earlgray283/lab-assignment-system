package repository

import (
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
)

type LabSurvey struct {
	UID       string
	Priority  int
	LabId     string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

const KindLabSurvey = "labSurvey"

func NewLabSurveyKey(uid string, priority int) *datastore.Key {
	return datastore.NameKey(KindLabSurvey, fmt.Sprintf("%s_%d", uid, priority), nil)
}

func NewLabSurvey(uid string, priority int, labId string) (*LabSurvey, *datastore.Key) {
	return &LabSurvey{
		UID:       uid,
		Priority:  priority,
		LabId:     labId,
		CreatedAt: time.Now(),
	}, NewLabSurveyKey(uid, priority)
}

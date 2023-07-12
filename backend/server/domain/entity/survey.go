package entity

import (
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
)

// 集計期間を定義する kind
type Survey struct {
	Year      int
	StartAt   time.Time
	EndAt     time.Time
	CreatedAt time.Time
	UpdatedAt *time.Time
}

const KindSurvey = "survey"

func NewSurveyKey(year int) *datastore.Key {
	return datastore.NameKey(KindSurvey, strconv.Itoa(year), nil)
}

func NewSurvey(year int, startAt, endAt, createdAt time.Time) (*Survey, *datastore.Key) {
	return &Survey{
		Year:      year,
		StartAt:   startAt,
		EndAt:     endAt,
		CreatedAt: createdAt,
	}, NewSurveyKey(year)
}

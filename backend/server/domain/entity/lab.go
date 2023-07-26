package entity

import (
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
)

type Lab struct {
	ID        string
	Name      string
	Lower     int
	Capacity  int
	Year      int
	UserGPAs  []*UserGPA // the length must be ApplicantCount
	Confirmed bool
	IsSpecial bool
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type UserGPA struct {
	UserKey *datastore.Key
	GPA     float64
}

const KindLab = "lab"
const LabsCapacity = 1024

func NewLabKey(labId string, year int) *datastore.Key {
	return datastore.NameKey(KindLab, fmt.Sprintf("%s_%d", labId, year), nil)
}

func NewLab(id string, name string, lower, capacity, year int, isSpecial bool, createdAt time.Time) (*Lab, *datastore.Key) {
	return &Lab{
		ID:        id,
		Name:      name,
		Lower:     lower,
		Capacity:  capacity,
		Year:      year,
		IsSpecial: isSpecial,
		CreatedAt: createdAt,
	}, NewLabKey(id, year)
}

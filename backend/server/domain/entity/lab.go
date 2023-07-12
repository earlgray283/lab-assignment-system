package entity

import (
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
)

type Lab struct {
	ID        string
	Name      string
	Capacity  int
	Year      int
	UserGPAs  []*UserGPA // the length must be ApplicantCount
	Confirmed bool
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

func NewLab(id string, name string, capacity, year int, createdAt time.Time) (*Lab, *datastore.Key) {
	return &Lab{
		ID:        id,
		Name:      name,
		Capacity:  capacity,
		Year:      year,
		CreatedAt: createdAt,
	}, NewLabKey(id, year)
}

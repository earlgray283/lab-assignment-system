package repository

import (
	"time"

	"cloud.google.com/go/datastore"
)

type Lab struct {
	ID        string
	Name      string
	Capacity  int
	CreatedAt time.Time
}

const KindLab = "lab"

func NewLabKey(labId string) *datastore.Key {
	return datastore.NameKey(KindLab, labId, nil)
}

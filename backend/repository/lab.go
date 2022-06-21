package repository

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
)

type Lab struct {
	ID           string
	Name         string
	Capacity     int
	CreatedAt    time.Time
}

const KindLab = "lab"
const LabsCapacity = 1024

func NewLabKey(labId string) *datastore.Key {
	return datastore.NameKey(KindLab, labId, nil)
}

// if len(labIds) == 0, fetch ALL labs
func FetchAllLabs(ctx context.Context, c *datastore.Client, labIds []string) ([]*Lab, bool, error) {
	if len(labIds) == 0 {
		repoLabs := make([]*Lab, 0, LabsCapacity)
		if _, err := c.GetAll(ctx, datastore.NewQuery(KindLab), &repoLabs); err != nil {
			return nil, false, err
		}
		return repoLabs, true, nil
	}

	repoLabs := make([]*Lab, len(labIds))
	repoKeys := make([]*datastore.Key, len(labIds))
	for i, labId := range labIds {
		repoKeys[i] = NewLabKey(labId)
	}
	if err := c.GetMulti(ctx, repoKeys, repoLabs); err != nil {
		if merr, ok := err.(datastore.MultiError); ok {
			for _, err := range merr {
				if err == datastore.ErrNoSuchEntity {
					return nil, false, nil
				}
			}
		}
		return nil, false, err
	}
	return repoLabs, true, nil
}

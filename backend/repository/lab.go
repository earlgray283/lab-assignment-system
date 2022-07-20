package repository

import (
	"context"
	"lab-assignment-system-backend/server/models"
	"time"

	"cloud.google.com/go/datastore"
	"golang.org/x/exp/slices"
)

type Lab struct {
	ID              string
	Name            string
	Capacity        int
	ConfirmedNumber int
	CreatedAt       time.Time
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

func CalculateLabGpa(c *datastore.Client) (map[string]*models.LabGpa, error) {
	ctx := context.Background()
	m := map[string]*models.LabGpa{}
	users := make([]*User, 0)

	if _, err := c.GetAll(ctx, datastore.NewQuery(KindUser), &users); err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Lab1 == nil || user.Lab2 == nil || user.Lab3 == nil {
			continue
		}
		if _, ok := m[*user.Lab1]; !ok {
			m[*user.Lab1] = models.NewLabGpa()
		}
		if _, ok := m[*user.Lab2]; !ok {
			m[*user.Lab2] = models.NewLabGpa()
		}
		if _, ok := m[*user.Lab3]; !ok {
			m[*user.Lab3] = models.NewLabGpa()
		}
		m[*user.Lab1].Gpas1 = append(m[*user.Lab1].Gpas1, user.Gpa)
		m[*user.Lab2].Gpas2 = append(m[*user.Lab2].Gpas2, user.Gpa)
		m[*user.Lab3].Gpas3 = append(m[*user.Lab3].Gpas3, user.Gpa)
	}

	cmpFunc := func(a, b float64) bool {
		return a > b
	}
	for _, labGpa := range m {
		slices.SortFunc(labGpa.Gpas1, cmpFunc)
		slices.SortFunc(labGpa.Gpas2, cmpFunc)
		slices.SortFunc(labGpa.Gpas3, cmpFunc)
	}

	return m, nil
}

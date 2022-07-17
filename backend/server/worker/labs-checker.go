package worker

import (
	"context"
	"lab-assignment-system-backend/repository"
	"lab-assignment-system-backend/server/models"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/datastore"
)

type LabMap struct {
	mp map[string]*models.LabGpa
	sync.Mutex
}

type LabsChecker struct {
	c        *datastore.Client
	labMap   *LabMap
	interval time.Duration
}

func NewLabsChecker(ctx context.Context, c *datastore.Client, interval time.Duration) (*LabsChecker, error) {
	keys, err := c.GetAll(ctx, datastore.NewQuery(repository.KindLab).KeysOnly(), nil)
	if err != nil {
		return nil, err
	}
	labGpaMp, err := repository.CalculateLabGpa(c)
	if err != nil {
		return nil, err
	}
	mp := map[string]*models.LabGpa{}
	for _, key := range keys {
		if labGpa, ok := labGpaMp[key.Name]; ok {
			mp[key.Name] = labGpa
		} else {
			mp[key.Name] = &models.LabGpa{
				Gpas1:     make([]float64, 0),
				Gpas2:     make([]float64, 0),
				Gpas3:     make([]float64, 0),
				UpdatedAt: time.Now(),
			}
		}
	}
	return &LabsChecker{
		c:        c,
		interval: interval,
		labMap:   &LabMap{mp: mp},
	}, nil
}

func (l *LabsChecker) GetLabGpa(labId string) *models.LabGpa {
	l.labMap.Lock()
	defer l.labMap.Unlock()
	return l.labMap.mp[labId]
}

func (l *LabsChecker) SingleRun() error {
	log.Println("Running Labs Checker...")

	l.labMap.Lock()
	defer l.labMap.Unlock()
	labMap, err := repository.CalculateLabGpa(l.c)
	if err != nil {
		return err
	}
	for labId, labGpa := range labMap {
		l.labMap.mp[labId] = labGpa
	}

	return nil
}

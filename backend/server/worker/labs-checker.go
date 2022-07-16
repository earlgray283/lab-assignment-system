package worker

import (
	"context"
	"lab-assignment-system-backend/lib"
	"lab-assignment-system-backend/repository"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/datastore"
)

type LabCountMap struct {
	mp map[string]*Priority
	sync.Mutex
}

func (l *LabCountMap) GetOrInsert(labId string, alt *Priority) *Priority {
	if l.mp[labId] == nil {
		l.mp[labId] = alt
	}
	return l.mp[labId]
}

type Priority struct {
	First, Second, Third int
}

type LabsChecker struct {
	c           *datastore.Client
	labCountMap *LabCountMap
	interval    time.Duration
}

func NewLabsChecker(c *datastore.Client, interval time.Duration) *LabsChecker {
	labCountMap := &LabCountMap{}
	return &LabsChecker{
		c:           c,
		interval:    interval,
		labCountMap: labCountMap,
	}
}

func (l *LabsChecker) GetLabCount(labId string) *Priority {
	l.labCountMap.Lock()
	defer l.labCountMap.Unlock()
	v := l.labCountMap.mp[labId]
	if v == nil {
		v = &Priority{}
	}
	return v
}

func (l *LabsChecker) GetLabCountMap() map[string]*Priority {
	l.labCountMap.Lock()
	defer l.labCountMap.Unlock()
	mp := map[string]*Priority{}
	for k, v := range l.labCountMap.mp {
		mp[k] = &Priority{
			First:  v.First,
			Second: v.Second,
			Third:  v.Third,
		}
	}
	return mp
}

func (l *LabsChecker) SingleRun() error {
	ctx := context.Background()
	log.Println("Running Labs Checker...")
	var users []*repository.User
	if _, err := l.c.GetAll(ctx, datastore.NewQuery(repository.KindUser), &users); err != nil {
		return err
	}

	l.labCountMap.Lock()
	defer l.labCountMap.Unlock()

	labCountMap := lib.Map[string, *Priority]{}
	for _, user := range users {
		labCountMap.GetOrInsert(*user.Lab1, &Priority{}).First++
		labCountMap.GetOrInsert(*user.Lab2, &Priority{}).Second++
		labCountMap.GetOrInsert(*user.Lab3, &Priority{}).Third++
	}
	l.labCountMap.mp = labCountMap

	return nil
}

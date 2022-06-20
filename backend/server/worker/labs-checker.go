package worker

import (
	"context"
	"fmt"
	"lab-assignment-system-backend/repository"
	"log"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/robfig/cron"
)

type LabsChecker struct {
	c        *datastore.Client
	interval time.Duration
}

func NewLabsChecker(c *datastore.Client, interval time.Duration) *LabsChecker {
	return &LabsChecker{
		c:        c,
		interval: interval,
	}
}

func (l *LabsChecker) Run() {
	c := cron.New()
	l.runnerFunc(c)()
	_ = c.AddFunc(fmt.Sprintf("@every %v", l.interval.String()), l.runnerFunc(c))
	c.Run()
}

func (l *LabsChecker) runnerFunc(c *cron.Cron) func() {
	return func() {
		ctx := context.Background()
		log.Println("Running Labs Checker...")
		var users []*repository.User
		if _, err := l.c.GetAll(ctx, datastore.NewQuery(repository.KindUser), &users); err != nil {
			log.Println(err)
			c.Stop()
			return
		}
		type Tuple struct {
			First, Second, Third int
		}
		labCountMap := map[string]*Tuple{}
		for _, user := range users {
			if _, ok := labCountMap[user.Lab1]; !ok {
				labCountMap[user.Lab1] = &Tuple{}
			}
			if _, ok := labCountMap[user.Lab2]; !ok {
				labCountMap[user.Lab2] = &Tuple{}
			}
			if _, ok := labCountMap[user.Lab3]; !ok {
				labCountMap[user.Lab3] = &Tuple{}
			}
			labCountMap[user.Lab1].First++
			labCountMap[user.Lab2].Second++
			labCountMap[user.Lab3].Third++
		}
		var labs []*repository.Lab
		if _, err := l.c.GetAll(ctx, datastore.NewQuery(repository.KindLab), &labs); err != nil {
			log.Println(err)
			c.Stop()
			return
		}
		for _, lab := range labs {
			if _, ok := labCountMap[lab.ID]; !ok {
				continue
			}
			if lab.FirstChoice != labCountMap[lab.ID].First {
				log.Printf("In lab \"%s\", first count is inconsistent(datastore: %v, labCountMap: %v)\n", lab.Name, lab.FirstChoice, labCountMap[lab.ID].First)
			}
			if lab.SecondChoice != labCountMap[lab.ID].Second {
				log.Printf("In lab \"%s\", second count is inconsistent(datastore: %v, labCountMap: %v)\n", lab.Name, lab.SecondChoice, labCountMap[lab.ID].Second)
			}
			if lab.ThirdChice != labCountMap[lab.ID].Third {
				log.Printf("In lab \"%s\", third count is inconsistent(datastore: %v, labCountMap: %v)\n", lab.Name, lab.ThirdChice, labCountMap[lab.ID].Third)
			}
		}
	}
}

package server

import (
	"context"
	"fmt"
	"lab-assignment-system-backend/repository"
	"log"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/robfig/cron"
	"golang.org/x/exp/slices"
)

type LabGpa struct {
	Gpas1     []float64
	Gpas2     []float64
	Gpas3     []float64
	UpdatedAt time.Time
}

type GpaWorker struct {
	c        *datastore.Client
	interval time.Duration
	m        map[string]*LabGpa
}

func NewGpaWorker(c *datastore.Client, interval time.Duration) *GpaWorker {
	return &GpaWorker{
		c:        c,
		interval: interval,
		m:        map[string]*LabGpa{},
	}
}

func (g *GpaWorker) Run() {
	c := cron.New()
	g.runnerFunc(c)()
	_ = c.AddFunc(fmt.Sprintf("@every %v", g.interval.String()), g.runnerFunc(c))
	c.Run()
}

func (g *GpaWorker) runnerFunc(c *cron.Cron) func() {
	return func() {
		log.Println("Running Gpa Worker...")
		ctx := context.Background()
		m := map[string]*LabGpa{}
		users := make([]*repository.User, 0)
		log.Println("Fetching All Users...")
		if _, err := g.c.GetAll(ctx, datastore.NewQuery(repository.KindUser), &users); err != nil {
			c.Stop()
			return
		}
		log.Println("Totaling Users Gpa...")
		for _, user := range users {
			if user.Gpa == nil {
				continue
			}
			if _, ok := m[user.Lab1]; !ok {
				m[user.Lab1] = &LabGpa{}
			}
			if _, ok := m[user.Lab2]; !ok {
				m[user.Lab2] = &LabGpa{}
			}
			if _, ok := m[user.Lab3]; !ok {
				m[user.Lab3] = &LabGpa{}
			}
			m[user.Lab1].Gpas1 = append(m[user.Lab1].Gpas1, *user.Gpa)
			m[user.Lab2].Gpas2 = append(m[user.Lab2].Gpas2, *user.Gpa)
			m[user.Lab3].Gpas3 = append(m[user.Lab3].Gpas3, *user.Gpa)
		}

		cmpFunc := func(a, b float64) bool {
			return a > b
		}
		for _, labGpa := range m {
			slices.SortFunc(labGpa.Gpas1, cmpFunc)
			slices.SortFunc(labGpa.Gpas2, cmpFunc)
			slices.SortFunc(labGpa.Gpas3, cmpFunc)
		}
		g.m = m
		log.Println("done")
	}
}

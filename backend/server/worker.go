package server

import (
	"context"
	"fmt"
	"lab-assignment-system-backend/repository"
	"lab-assignment-system-backend/server/models"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/robfig/cron"
	"golang.org/x/exp/slices"
)

type LabGpaStorage struct {
	m    map[string]*models.LabGpa
	gpas []float64
	sync.RWMutex
}

type GpaWorker struct {
	c        *datastore.Client
	interval time.Duration
	m        *LabGpaStorage
}

func NewGpaWorker(c *datastore.Client, interval time.Duration) *GpaWorker {
	return &GpaWorker{
		c:        c,
		interval: interval,
		m:        &LabGpaStorage{m: map[string]*models.LabGpa{}},
	}
}

func (g *GpaWorker) Run() {
	c := cron.New()
	g.runnerFunc(c)()
	_ = c.AddFunc(fmt.Sprintf("@every %v", g.interval.String()), g.runnerFunc(c))
	c.Run()
}

func (g *GpaWorker) GetGpas() []float64 {
	return g.m.gpas
}

func (g *GpaWorker) Get(labId string) *models.LabGpa {
	g.m.Lock()
	defer g.m.Unlock()
	return g.m.m[labId]
}

func (g *GpaWorker) runnerFunc(c *cron.Cron) func() {
	return func() {
		log.Println("Running Gpa Worker...")
		ctx := context.Background()
		m := map[string]*models.LabGpa{}
		users := make([]*repository.User, 0)
		log.Println("Fetching All Users...")
		if _, err := g.c.GetAll(ctx, datastore.NewQuery(repository.KindUser), &users); err != nil {
			c.Stop()
			return
		}
		log.Println("Totaling Users Gpa...")
		gpas := make([]float64, len(users))
		for i, user := range users {
			if user.Gpa == nil {
				continue
			}
			gpas[i] = *user.Gpa
			if _, ok := m[user.Lab1]; !ok {
				m[user.Lab1] = &models.LabGpa{}
			}
			if _, ok := m[user.Lab2]; !ok {
				m[user.Lab2] = &models.LabGpa{}
			}
			if _, ok := m[user.Lab3]; !ok {
				m[user.Lab3] = &models.LabGpa{}
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
		g.m.Lock()
		g.m.m = m
		g.m.gpas = gpas
		g.m.Unlock()
		log.Println("done")
	}
}

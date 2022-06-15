package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

type Lab struct {
	ID           string
	Name         string
	Capacity     int
	FirstChoice  int
	SecondChoice int
	ThirdChice   int
	CreatedAt    time.Time
}

const KindLab = "lab"

func NewLabKey(labId string) *datastore.Key {
	return datastore.NameKey(KindLab, labId, nil)
}

const ProjectId = "lab-assignment-system-project"

func main() {
	dc, err := datastore.NewClient(context.Background(), ProjectId)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open("labs.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	mutations := make([]*datastore.Mutation, 0, 1024)
	for sc.Scan() {
		tokens := strings.Split(sc.Text(), ",")
		capacity, _ := strconv.Atoi(tokens[2])
		lab := &Lab{
			ID:        tokens[1],
			Name:      tokens[0],
			Capacity:  capacity,
			CreatedAt: time.Now(),
		}
		mutations = append(mutations, datastore.NewInsert(NewLabKey(lab.ID), lab))
	}
	if _, err := dc.Mutate(context.Background(), mutations...); err != nil {
		log.Fatal(err)
	}
}

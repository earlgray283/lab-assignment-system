package main

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

const ProjectId = "lab-assignment-system-project"

type Lab struct {
	ID              string
	Name            string
	Capacity        int
	ConfirmedNumber int
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

const KindLab = "lab"

func NewLabKey(id string) *datastore.Key {
	return datastore.NameKey(KindLab, id, nil)
}

func main() {
	d, err := datastore.NewClient(context.Background(), ProjectId, option.WithCredentialsFile("../../backend/credentials.json"))
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	labs := make([]Lab, 0)
	labKeys := make([]*datastore.Key, 0)
	csvf := csv.NewReader(f)
	for {
		cols, err := csvf.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		capacity, _ := strconv.Atoi(cols[1])
		now := time.Now()
		labs = append(labs, Lab{
			ID:        cols[2],
			Name:      cols[0],
			Capacity:  capacity,
			UpdatedAt: &now,
		})
		labKeys = append(labKeys, NewLabKey(cols[2]))
	}

	if _, err := d.PutMulti(context.Background(), labKeys, labs); err != nil {
		log.Fatal(err)
	}
}

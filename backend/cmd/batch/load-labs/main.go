package main

import (
	"context"
	"encoding/csv"
	"flag"
	"io"
	"lab-assignment-system-backend/server/domain/entity"
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
)

var year = flag.Int("year", time.Now().Year(), "year")

func main() {
	flag.Parse()

	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Fatal("please set GCP_PROJECT_ID")
	}
	csvPath := flag.Arg(0)
	if csvPath == "" {
		flag.Usage()
		return
	}

	dc, err := datastore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatal(err)
	}
	defer dc.Close()

	f, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	mutations := make([]*datastore.Mutation, 0)
	for {
		records, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		capacity, err := strconv.Atoi(records[2])
		if err != nil {
			log.Fatal(err)
		}
		isSpecial, err := strconv.ParseBool(records[3])
		if err != nil {
			log.Fatal(err)
		}
		lab := &entity.Lab{
			ID:        records[0],
			Name:      records[1],
			Capacity:  capacity,
			Year:      *year,
			IsSpecial: isSpecial,
			CreatedAt: time.Now(),
		}
		mutations = append(mutations, datastore.NewUpsert(entity.NewLabKey(lab.ID, *year), lab))
	}
	if _, err := dc.Mutate(context.Background(), mutations...); err != nil {
		log.Fatal(err)
	}
}

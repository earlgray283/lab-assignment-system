package main

import (
	"context"
	"flag"
	"fmt"
	"lab-assignment-system-backend/server/domain/entity"
	"log"
	"os"
	"time"

	"cloud.google.com/go/datastore"
)

var (
	year        = *flag.Int("year", time.Now().Year(), "year")
	startAtText = *flag.String("startAt", "", "format: 2006-01-02T15:04:05")
	endAtText   = *flag.String("endAt", "", "format: 2006-01-02T15:04:05")
)

func main() {
	flag.Parse()
	startAt, err := time.Parse("2006-01-02T15:04:05", startAtText)
	if err != nil {
		log.Fatal(err)
	}
	endAt, err := time.Parse("2006-01-02T15:04:05", endAtText)
	if err != nil {
		log.Fatal(err)
	}

	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Fatal("please set GCP_PROJECT_ID")
	}
	dc, err := datastore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatal(err)
	}
	defer dc.Close()

	survey, key := entity.NewSurvey(year, startAt, endAt, time.Now())
	if _, err := dc.Put(context.Background(), key, survey); err != nil {
		log.Fatal(err)
	}
	fmt.Println(*survey)
}

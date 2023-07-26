package main

import (
	"context"
	"fmt"
	"io"
	"lab-assignment-system-backend/server/usecases"
	"log"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
)

const ProjectId = "lab-assignment-system-project"

func main() {
	ctx := context.Background()
	dsClient, err := datastore.NewClient(context.Background(), ProjectId)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer dsClient.Close()
	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer gcsClient.Close()

	labCSV, userCSV, err := usecases.NewAdminInteractor(dsClient, log.Default()).FinalDecisionDryRun(ctx, 2023)
	if err != nil {
		log.Fatal(err)
		return
	}

	bucket := gcsClient.Bucket("lab-result")
	now := time.Now().Format("2006-01-02T1504")

	labw := bucket.Object(fmt.Sprintf("%s/lab.csv", now)).NewWriter(ctx)
	if _, err := io.Copy(labw, labCSV); err != nil {
		log.Fatal(err)
		return
	}
	if err := labw.Close(); err != nil {
		log.Fatal(err)
		return
	}

	userw := bucket.Object(fmt.Sprintf("%s/user.csv", now)).NewWriter(ctx)
	if _, err := io.Copy(userw, userCSV); err != nil {
		log.Fatal(err)
		return
	}
	if err := userw.Close(); err != nil {
		log.Fatal(err)
		return
	}
}

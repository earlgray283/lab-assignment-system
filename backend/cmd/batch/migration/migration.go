package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/datastore"
)

func main() {
	ctx := context.Background()
	dsClient, err := datastore.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"))
	if err != nil {
		log.Fatal(err)
	}
	defer dsClient.Close()

	// write here
}

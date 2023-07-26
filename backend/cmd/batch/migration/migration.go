package main

import (
	"context"
	"lab-assignment-system-backend/server/domain/entity"
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

	oldLabs := make([]*entity.Lab, 0)
	keys, err := dsClient.GetAll(ctx, datastore.NewQuery(entity.KindLab), &oldLabs)
	if err != nil {
		log.Fatal(err)
	}
	for _, oldLab := range oldLabs {
		if oldLab.Capacity == 5 {
			oldLab.Lower = 4
		} else if oldLab.Capacity == 4 {
			oldLab.Lower = 4
		}
	}
	if _, err := dsClient.PutMulti(ctx, keys, oldLabs); err != nil {
		log.Fatal(err)
	}
}

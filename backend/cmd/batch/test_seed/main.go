package main

import (
	"context"
	"lab-assignment-system-backend/server/domain/entity"
	"lab-assignment-system-backend/server/lib"
	"log"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/samber/lo"
)

func main() {
	ctx := context.Background()
	dsClient, err := datastore.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"))
	if err != nil {
		log.Fatal(err)
	}
	defer dsClient.Close()

	users := []struct {
		uid     string
		wishLab *string
		gpa     float64
		year    int
	}{
		{"1", lo.ToPtr("ohkilab"), 2.9, 2023},
		{"2", lo.ToPtr("ohkilab"), 3.0, 2023},
		{"3", lo.ToPtr("ohkilab"), 3.1, 2023},
	}
	labs := []struct {
		id       string
		capacity int
		year     int
	}{
		{"ohkilab", 2, 2023},
		{"uhkilab", 2, 2023},
		{"ahkilab", 2, 2023},
	}

	survey, surveyKey := entity.NewSurvey(2023, lib.YMD(2023, 1, 1), lib.YMD(2023, 12, 31), time.Now())
	if _, err := dsClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		mutations := make([]*datastore.Mutation, 0)
		for _, user := range users {
			newUser, key := entity.NewUser(user.uid, user.gpa, user.year, entity.RoleAudience, time.Now())
			newUser.WishLab = user.wishLab
			mutations = append(mutations, datastore.NewUpsert(key, newUser))
		}
		for _, lab := range labs {
			newLab, key := entity.NewLab(lab.id, lab.id, lab.capacity, lab.year, time.Now())
			for _, user := range users {
				if user.wishLab == nil {
					continue
				}
				if *user.wishLab != lab.id {
					continue
				}
				newLab.UserGPAs = append(newLab.UserGPAs, &entity.UserGPA{
					UserKey: entity.NewUserKey(user.uid),
					GPA:     user.gpa,
				})
			}
			mutations = append(mutations, datastore.NewUpsert(key, newLab))
		}
		mutations = append(mutations, datastore.NewUpsert(surveyKey, survey))
		_, err := tx.Mutate(mutations...)
		return err
	}); err != nil {
		log.Fatal(err)
	}
}

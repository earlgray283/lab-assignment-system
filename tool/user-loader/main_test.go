package main

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

func TestFetchAllUsers(t *testing.T) {
	d, _ := datastore.NewClient(context.Background(), ProjectId, option.WithCredentialsFile("../../backend/credentials.json"))
	users, _, err := fetchAllUsers(context.Background(), d)
	if err != nil {
		t.Fatal(err)
	}
	for _, user := range users {
		fmt.Println(user.UID, user.Gpa, user.Lab1, user.Lab2, user.Lab3)
	}
}

func TestLoadUsersFromCsv(t *testing.T) {
	d, _ := datastore.NewClient(context.Background(), ProjectId, option.WithCredentialsFile("../../backend/credentials.json"))
	users, userKeys, err := fetchAllUsers(context.Background(), d)
	if err != nil {
		t.Fatal(err)
	}
	userKeyNames := mapSlice(userKeys, func(k *datastore.Key) string {
		return k.Name
	})
	userMap := NewMapFrom2Slices(userKeyNames, users)

	newUsers, _, err := loadUsersFromCsv("data.csv", userMap)
	if err != nil {
		t.Fatal(err)
	}
	hasLabCount := 0
	for _, user := range newUsers {
		if user.Lab1 != nil && user.Lab2 != nil && user.Lab3 != nil {
			hasLabCount++
		}
		if userMap[user.UID].Lab1 != user.Lab1 || userMap[user.UID].Lab2 != user.Lab2 || userMap[user.UID].Lab3 != user.Lab3 {
			t.Fatal("labs not equal")
		}
	}
	if hasLabCount == 0 {
		t.Fatal("lab count must not be 0")
	}
	t.Log(hasLabCount)
}

package main

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

const ProjectId = "lab-assignment-system-project"

type User struct {
	UID  string  `json:"uid,omitempty"`
	Gpa  float64 `json:"gpa"`
	Lab1 *string `json:"lab1,omitempty"`
	Lab2 *string `json:"lab2,omitempty"`
	Lab3 *string `json:"lab3,omitempty"`
}

const KindUser = "user"

func NewUserKey(uid string) *datastore.Key {
	return datastore.NameKey(KindUser, uid, nil)
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

	users := make([]User, 0)
	userKeys := make([]*datastore.Key, 0)
	r := csv.NewReader(f)
	for {
		cols, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		gpa, _ := strconv.ParseFloat(cols[0], 64)
		users = append(users, User{
			UID: cols[1],
			Gpa: gpa,
		})
		userKeys = append(userKeys, NewUserKey(cols[1]))
	}
	if _, err := d.PutMulti(context.Background(), userKeys, users); err != nil {
		log.Fatal(err)
	}
}

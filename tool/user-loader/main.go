package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

const ProjectId = "lab-assignment-system-project"

type User struct {
	UID          string  `json:"uid,omitempty"`
	Gpa          float64 `json:"gpa"`
	Lab1         *string `json:"lab1,omitempty"`
	Lab2         *string `json:"lab2,omitempty"`
	Lab3         *string `json:"lab3,omitempty"`
	ConfirmedLab string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

const KindUser = "user"

func NewUserKey(uid string) *datastore.Key {
	return datastore.NameKey(KindUser, uid, nil)
}

type Map[K comparable, V any] map[K]V

func NewMapFrom2Slices[K comparable, V any](keys []K, values []V) map[K]V {
	if len(keys) != len(values) {
		panic("")
	}
	mp := make(map[K]V, len(keys))
	for i := range keys {
		mp[keys[i]] = values[i]
	}
	return mp
}

func mapSlice[T, U any](a []T, f func(t T) U) []U {
	newa := make([]U, len(a))
	for i, elem := range a {
		newa[i] = f(elem)
	}
	return newa
}

func main() {
	d, err := datastore.NewClient(context.Background(), ProjectId, option.WithCredentialsFile("../../backend/credentials.json"))
	if err != nil {
		log.Fatal(err)
	}

	users, userKeys, err := fetchAllUsers(context.Background(), d)
	if err != nil {
		log.Fatal(err)
	}
	userKeyNames := mapSlice(userKeys, func(k *datastore.Key) string {
		return k.Name
	})
	userMap := NewMapFrom2Slices(userKeyNames, users)

	newUsers, newUserKeys, err := loadUsersFromCsv("data.csv", userMap)
	if err != nil {
		log.Fatal(err)
	}
	_ = newUserKeys
	for _, user := range newUsers {
		fmt.Println(user.Lab1, user.Lab2, user.Lab3)
	}

	if _, err := d.PutMulti(context.Background(), newUserKeys, newUsers); err != nil {
		log.Fatal(err)
	}
}

func fetchAllUsers(ctx context.Context, d *datastore.Client) ([]*User, []*datastore.Key, error) {
	var users []*User
	keys, err := d.GetAll(ctx, datastore.NewQuery(KindUser), &users)
	if err != nil {
		return nil, nil, err
	}
	return users, keys, nil
}

func loadUsersFromCsv(csvName string, oldUserMap map[string]*User) ([]*User, []*datastore.Key, error) {
	f, err := os.Open(csvName)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	newUsers := make([]*User, 0)
	newUserKeys := make([]*datastore.Key, 0)
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
		newUser := oldUserMap[NewUserKey(cols[1]).Name]
		if newUser == nil {
			newUser = &User{
				UID:       cols[1],
				Gpa:       gpa,
				CreatedAt: time.Now(),
			}
		}
		newUserKeys = append(newUserKeys, NewUserKey(cols[1]))
		newUsers = append(newUsers, newUser)
	}
	return newUsers, newUserKeys, nil
}

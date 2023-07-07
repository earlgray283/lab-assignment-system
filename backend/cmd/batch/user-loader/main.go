package main

import (
	"context"
	"encoding/csv"
	"io"
	"lab-assignment-system-backend/server/domain/entity"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

const ProjectId = "lab-assignment-system-project"

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
	dc, err := datastore.NewClient(context.Background(), ProjectId)
	if err != nil {
		log.Fatal(err)
	}
	defer dc.Close()

	users, userKeys, err := fetchAllUsers(context.Background(), dc)
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

	if _, err := dc.PutMulti(context.Background(), newUserKeys, newUsers); err != nil {
		log.Fatal(err)
	}
}

func fetchAllUsers(ctx context.Context, d *datastore.Client) ([]*entity.User, []*datastore.Key, error) {
	var users []*entity.User
	keys, err := d.GetAll(ctx, datastore.NewQuery(entity.KindUser), &users)
	if err != nil {
		return nil, nil, err
	}
	return users, keys, nil
}

func loadUsersFromCsv(csvName string, oldUserMap map[string]*entity.User) ([]*entity.User, []*datastore.Key, error) {
	f, err := os.Open(csvName)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	newUsers := make([]*entity.User, 0)
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
		uid := strings.Trim(cols[0], "\ufeff")
		gpa, err := strconv.ParseFloat(uid, 64)
		if err != nil {
			return nil, nil, err
		}
		newUser := oldUserMap[entity.NewUserKey(cols[1]).Name]
		if newUser == nil {
			newUser = &entity.User{
				UID:       cols[1],
				Gpa:       gpa,
				CreatedAt: time.Now(),
			}
		}
		now := time.Now()
		newUser.Gpa = gpa
		newUser.UpdatedAt = &now
		newUserKeys = append(newUserKeys, entity.NewUserKey(cols[1]))
		newUsers = append(newUsers, newUser)
	}
	return newUsers, newUserKeys, nil
}

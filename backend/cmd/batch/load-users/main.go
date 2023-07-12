package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"lab-assignment-system-backend/server/domain/entity"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

var year = flag.Int("year", time.Now().Year(), "year")

func main() {
	flag.Parse()
	ctx := context.Background()

	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Fatal("please set GCP_PROJECT_ID")
	}
	csvPath := flag.Arg(0)
	if csvPath == "" {
		log.Fatal("please specify csv location")
	}

	dc, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
	defer dc.Close()

	// csv から load
	newUsers, newUserKeys, err := loadUsersFromCsv(csvPath)
	if err != nil {
		log.Fatal(err)
	}

	// 更新(upsert)
	if _, err := dc.PutMulti(ctx, newUserKeys, newUsers); err != nil {
		log.Fatal(err)
	}
}

func loadUsersFromCsv(csvName string) ([]*entity.User, []*datastore.Key, error) {
	f, err := os.Open(csvName)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	users := make([]*entity.User, 0)
	userKeys := make([]*datastore.Key, 0)
	r := csv.NewReader(f)
	line := 1
	for {
		cols, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		user, key, err := parseLine(cols)
		if err != nil {
			return nil, nil, fmt.Errorf("line %d: %w", line, err)
		}
		users = append(users, user)
		userKeys = append(userKeys, key)

		line++
	}
	return users, userKeys, nil
}

func parseLine(cols []string) (*entity.User, *datastore.Key, error) {
	uid := strings.TrimSpace(cols[0])
	gpa, err := strconv.ParseFloat(cols[1], 64)
	if err != nil {
		return nil, nil, err
	}
	role, ok := entity.RoleByValue[cols[2]]
	if !ok {
		return nil, nil, errors.New("no such role")
	}
	user, key := entity.NewUser(uid, gpa, *year, role, time.Now())
	return user, key, nil
}

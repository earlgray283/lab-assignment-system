package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"lab-assignment-system-backend/server/domain/entity"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

const ProjectId = "lab-assignment-system-project"

var (
	year = flag.Int("year", time.Now().Year(), "year")
)

func main() {
	flag.Parse()

	csvPath := flag.Arg(0)
	if csvPath == "" {
		fmt.Println("please specify csv location")
	}

	dc, err := datastore.NewClient(context.Background(), ProjectId)
	if err != nil {
		log.Fatal(err)
	}
	defer dc.Close()

	f, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	mutations := make([]*datastore.Mutation, 0)
	for sc.Scan() {
		tokens := strings.Split(sc.Text(), ",")
		capacity, _ := strconv.Atoi(tokens[2])
		lab := &entity.Lab{
			ID:        tokens[1],
			Name:      tokens[0],
			Capacity:  capacity,
			CreatedAt: time.Now(),
		}
		mutations = append(mutations, datastore.NewInsert(entity.NewLabKey(lab.ID, *year), lab))
	}
	if _, err := dc.Mutate(context.Background(), mutations...); err != nil {
		log.Fatal(err)
	}
}

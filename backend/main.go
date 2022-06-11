package main

import (
	"context"
	"lab-assignment-system-backend/server"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
)

const ProjectId = "lab-assignment-system-project"

func init() {
	_ = godotenv.Load(".env")
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Fatal("env GOOGLE_APPLICATION_CREDENTIALS must be set")
	}
}

func main() {
	frontendUrl := os.Getenv("FRONT_URL")
	if frontendUrl == "" {
		log.Fatal("environmental value FRONT_URL must be set")
	}
	dc, err := datastore.NewClient(context.Background(), ProjectId)
	if err != nil {
		log.Fatal(err)
	}
	fa, err := firebase.NewApp(context.Background(), &firebase.Config{ProjectID: ProjectId})
	if err != nil {
		log.Fatal(err)
	}
	auth, err := fa.Auth(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	srv := server.New(dc, auth, frontendUrl)
	if err := srv.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"fmt"
	"lab-assignment-system-backend/server"
	"log"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
)

const ProjectId = "lab-assignment-system-project"

func init() {
	_ = godotenv.Load(".env")

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	time.Local = jst
}

func main() {
	frontendUrl := os.Getenv("FRONTEND_URL")
	if frontendUrl == "" {
		log.Fatal("environmental value FRONTEND_URL must be set")
	}
	gakujoUrl := os.Getenv("GAKUJO_URL")
	if gakujoUrl == "" {
		log.Fatal("environmental value GAKUJO_URL must be set")
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

	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	srv := server.New(dc, auth, frontendUrl, gakujoUrl)
	if err := srv.Run(fmt.Sprintf(":%v", port)); err != nil {
		log.Fatal(err)
	}
}

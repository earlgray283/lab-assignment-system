package main

import (
	"context"
	"fmt"
	"lab-assignment-system-backend/lib"
	"lab-assignment-system-backend/server"
	"log"
	"os"
	"time"

	_ "time/tzdata"

	"cloud.google.com/go/datastore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
)

const ProjectId = "lab-assignment-system-project"

var (
	frontendUrl    string
	gakujoUrl      string
	senderEmail    string
	senderPassword string
	senderSmtp     string
)

func getEnvOrFatal(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatal("environmental value " + key + " must be set")
	}
	return value
}

func init() {
	_ = godotenv.Load(".env")

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	time.Local = jst

	frontendUrl = getEnvOrFatal("FRONTEND_URL")
	gakujoUrl = getEnvOrFatal("GAKUJO_URL")
	senderEmail = getEnvOrFatal("SENDER_EMAIL")
	senderPassword = getEnvOrFatal("SENDER_PASSWORD")
	senderSmtp = getEnvOrFatal("SENDER_SMTP")
}

func main() {
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
	smtpCli := lib.NewSmtpClient(senderEmail, senderPassword, senderSmtp, "587")

	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	srv := server.New(dc, auth, smtpCli, []string{frontendUrl, gakujoUrl})
	if err := srv.Run(fmt.Sprintf(":%v", port)); err != nil {
		log.Fatal(err)
	}
}

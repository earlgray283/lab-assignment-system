package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

const ProjectId = "lab-assignment-system-project"

var uid string = "secret"

func main() {
	ctx := context.Background()
	fa, err := firebase.NewApp(context.Background(), &firebase.Config{ProjectID: ProjectId}, option.WithCredentialsFile("../../credentials.json"))
	if err != nil {
		log.Fatal(err)
	}
	auth, err := fa.Auth(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	customClaims := map[string]interface{}{"admin": true}
	if err := auth.SetCustomUserClaims(ctx, uid, customClaims); err != nil {
		log.Fatal(err)
	}
}

package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/datastore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

const ProjectId = "lab-assignment-system-project"

const gradesJsonText = `<secret>`
const registerToken = `<secret>` // TODO: 自動生成

func TestPostGrades(t *testing.T) {
	go launchServer()

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/grades", strings.NewReader(gradesJsonText))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Register-Token", registerToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("not 200 OK")
	}
}

func TestDeleteUser(t *testing.T) {
	go launchServer()
	uid := "<secret>"
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/users/"+uid, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Register-Token", registerToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("not 200 OK")
	}
}

func init() {
	_ = godotenv.Load("../.env")
}

func launchServer() {
	frontendUrl := os.Getenv("FRONTEND_URL")
	if frontendUrl == "" {
		log.Fatal("environmental value FRONTEND_URL must be set")
	}
	gakujoUrl := os.Getenv("GAKUJO_URL")
	if gakujoUrl == "" {
		log.Fatal("environmental value GAKUJO_URL must be set")
	}
	dc, err := datastore.NewClient(context.Background(), ProjectId, option.WithCredentialsFile("../credentials.json"))
	if err != nil {
		log.Fatal(err)
	}
	fa, err := firebase.NewApp(context.Background(), &firebase.Config{ProjectID: ProjectId}, option.WithCredentialsFile("../credentials.json"))
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
	srv := New(dc, auth, []string{frontendUrl, gakujoUrl})
	if err := srv.Run(fmt.Sprintf(":%v", port)); err != nil {
		log.Fatal(err)
	}
}

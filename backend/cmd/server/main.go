package main

import (
	"context"
	"fmt"
	"lab-assignment-system-backend/server/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "time/tzdata"

	"cloud.google.com/go/datastore"
	"github.com/joho/godotenv"
)

const ProjectId = "lab-assignment-system-project"

var (
	frontendUrl string
	gakujoUrl   string
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
}

func main() {
	dc, err := datastore.NewClient(context.Background(), ProjectId)
	if err != nil {
		log.Fatal(err)
	}
	defer dc.Close()

	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	addr := fmt.Sprintf(":%v", port)

	srv := api.NewServer(dc, addr, []string{frontendUrl, gakujoUrl})
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	<-ctx.Done()
	log.Println("timeout of 5 seconds.")
	log.Println("Server exiting")
}

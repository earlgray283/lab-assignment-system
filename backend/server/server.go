package server

import (
	"context"
	"lab-assignment-system-backend/server/worker"
	"log"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	r           *gin.Engine
	logger      *log.Logger
	dc          *datastore.Client
	labsChecker *worker.LabsChecker
}

func NewCorsConfig(allowOrigins []string) *cors.Config {
	config := cors.DefaultConfig()
	config.AllowOrigins = allowOrigins
	config.AllowCredentials = true
	return &config
}

func New(dc *datastore.Client, allowOrigins []string) (*Server, error) {
	r := gin.Default()
	corsConfig := NewCorsConfig(append([]string{"http://localhost:3000"}, allowOrigins...))
	r.Use(cors.New(*corsConfig))
	logger := log.Default()
	gin.DefaultWriter = logger.Writer()
	labsWorker, err := worker.NewLabsChecker(context.Background(), dc, time.Hour)
	if err != nil {
		return nil, err
	}

	srv := &Server{r, logger, dc, labsWorker}

	srv.AuthRouter()
	srv.LabsRouter()
	srv.UserRouter()
	srv.GradesRouter()

	return srv, nil
}

func (srv *Server) Run(addr ...string) error {
	errc := make(chan error)
	if err := srv.labsChecker.SingleRun(); err != nil {
		return err
	}
	go func() {
		err := srv.r.Run(addr...)
		if err != nil {
			errc <- err
		}
	}()
	return <-errc
}

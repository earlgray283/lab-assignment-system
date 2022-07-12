package server

import (
	"lab-assignment-system-backend/repository"
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
	labGpa      *repository.LabGpa
}

const ExcludeLowerPoint = 60

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
	labsWorker := worker.NewLabsChecker(dc, time.Hour)
	labGpa, err := repository.CalculateLabGpa(dc)
	if err != nil {
		return nil, err
	}
	srv := &Server{r, logger, dc, labsWorker, labGpa}

	srv.AuthRouter()
	srv.LabsRouter()
	srv.UserRouter()

	return srv, nil
}

func (srv *Server) Run(addr ...string) error {
	errc := make(chan error)
	go func() {
		srv.labsChecker.Run()
	}()
	go func() {
		err := srv.r.Run(addr...)
		if err != nil {
			errc <- err
		}
	}()
	return <-errc
}

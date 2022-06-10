package server

import (
	"log"

	"cloud.google.com/go/datastore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

type Server struct {
	r           *gin.Engine
	logger      *log.Logger
	dc          *datastore.Client
	auth        *auth.Client
	frontendUrl string
}

func New(dc *datastore.Client, auth *auth.Client, frontendUrl string) *Server {
	r := gin.Default()
	logger := log.Default()
	gin.DefaultWriter = logger.Writer()
	srv := &Server{r, logger, dc, auth, frontendUrl}

	srv.GradesRouter()

	return srv
}

func (srv *Server) Run(addr ...string) error {
	return srv.Run(addr...)
}

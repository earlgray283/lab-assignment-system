package server

import (
	"log"
	"time"

	"cloud.google.com/go/datastore"
	"firebase.google.com/go/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	r           *gin.Engine
	logger      *log.Logger
	dc          *datastore.Client
	auth        *auth.Client
	frontendUrl string
	gpaWorker   *GpaWorker
}

const ExcludeLowerPoint = 60

func NewCorsConfig(allowOrigins []string) *cors.Config {
	config := cors.DefaultConfig()
	config.AllowOrigins = allowOrigins
	config.AllowHeaders = append(config.AllowHeaders, "register-token")
	config.AllowCredentials = true
	return &config
}

func New(dc *datastore.Client, auth *auth.Client, frontendUrl, gakujoUrl string) *Server {
	r := gin.Default()
	corsConfig := NewCorsConfig([]string{
		"http://localhost:3000",
		gakujoUrl,
		frontendUrl,
	})
	r.Use(cors.New(*corsConfig))
	logger := log.Default()
	gin.DefaultWriter = logger.Writer()
	gpaWorker := NewGpaWorker(dc, time.Hour)
	srv := &Server{r, logger, dc, auth, frontendUrl, gpaWorker}

	srv.GradesRouter()
	srv.AuthRouter()
	srv.LabsRouter()
	srv.UserRouter()

	return srv
}

func (srv *Server) Run(addr ...string) error {
	errc := make(chan error)
	go func() {
		srv.gpaWorker.Run()
	}()
	go func() {
		err := srv.r.Run(addr...)
		if err != nil {
			errc <- err
		}
	}()
	return <-errc
}

func (srv *Server) GetAuthToken(c *gin.Context) (*auth.Token, error) {
	sessionCookie, err := c.Cookie("session")
	if err != nil {
		return nil, err
	}
	authToken, err := srv.auth.VerifySessionCookieAndCheckRevoked(c.Request.Context(), sessionCookie)
	if err != nil {
		return nil, err
	}
	return authToken, nil
}

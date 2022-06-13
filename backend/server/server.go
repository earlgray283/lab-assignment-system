package server

import (
	"log"

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
	srv := &Server{r, logger, dc, auth, frontendUrl}

	srv.GradesRouter()
	srv.AuthRouter()

	return srv
}

func (srv *Server) Run(addr ...string) error {
	return srv.r.Run(addr...)
}

func (srv *Server) GetAuthToken(c *gin.Context) (*auth.Token, error) {
	sessionCookie, err := c.Cookie("session")
	if err != nil {
		return nil, err
	}
	authToken, err := srv.auth.VerifySessionCookie(c.Request.Context(), sessionCookie)
	if err != nil {
		return nil, err
	}
	return authToken, nil
}

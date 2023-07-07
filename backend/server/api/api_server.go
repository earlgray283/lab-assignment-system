package api

import (
	"lab-assignment-system-backend/server/api/middleware"
	i_http "lab-assignment-system-backend/server/interfaces/http"
	"lab-assignment-system-backend/server/usecases"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

func NewServer(dsClient *datastore.Client, addr string, allowOrigins []string) *http.Server {
	r := gin.Default()
	r.Use(middleware.Cors(allowOrigins))

	logger := log.Default()
	authInteractor := usecases.NewAuthInteractor(dsClient, logger)
	authController := i_http.NewAuthController(authInteractor)
	labsInteractor := usecases.NewLabsInteractor(dsClient, logger)
	labsController := i_http.NewLabsController(labsInteractor)
	usersInteractor := usecases.NewUsersInteractor(dsClient, logger)
	usersController := i_http.NewUsersController(usersInteractor)

	r.POST("/auth/signin", authController.Login)
	r.POST("/auth/signout", authController.Logout, middleware.Authentication(dsClient))
	r.POST("/labs", labsController.ListLabs)
	// TODO: GET /labs/csv
	r.PUT("/users/lab", usersController.UpdateUser, middleware.Authentication(dsClient))

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

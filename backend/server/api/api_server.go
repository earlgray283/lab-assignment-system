package api

import (
	"lab-assignment-system-backend/server/api/middleware"
	i_http "lab-assignment-system-backend/server/interfaces/http"
	"lab-assignment-system-backend/server/usecases"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer(dsClient *datastore.Client, addr string, corsConfig *cors.Config) *http.Server {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(cors.New(*corsConfig))

	logger := log.Default()
	authInteractor := usecases.NewAuthInteractor(dsClient, logger)
	authController := i_http.NewAuthController(authInteractor)
	labsInteractor := usecases.NewLabsInteractor(dsClient, logger)
	labsController := i_http.NewLabsController(labsInteractor)
	usersInteractor := usecases.NewUsersInteractor(dsClient, logger)
	usersController := i_http.NewUsersController(usersInteractor)
	gradesInteractor := usecases.NewGradesInteractor(dsClient, logger)
	gradesController := i_http.NewGradesController(gradesInteractor)
	adminInteractor := usecases.NewAdminInteractor(dsClient, logger)
	adminController := i_http.NewAdminController(adminInteractor)

	// audience APIs
	r.POST("/auth/signin", authController.Login)
	r.POST("/auth/signout", middleware.Authentication(dsClient), authController.Logout)
	r.GET("/labs", labsController.ListLabs)
	r.PUT("/users/lab", middleware.Authentication(dsClient), usersController.UpdateUser)
	r.GET("/users/me", middleware.Authentication(dsClient), usersController.GetUserMe)
	r.GET("/grades", middleware.Authentication(dsClient), gradesController.ListGrades)

	// admin APIs
	r.POST("/admin/final-decision", middleware.Authentication(dsClient), adminController.FinalDecision)
	r.GET("/admin/csv", middleware.Authentication(dsClient), adminController.GetCSV)
	r.POST("/admin/create-users", middleware.Authentication(dsClient), adminController.CreateUsers)

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

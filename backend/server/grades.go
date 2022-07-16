package server

import (
	"lab-assignment-system-backend/repository"
	"lab-assignment-system-backend/server/models"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

func (srv *Server) GradesRouter() {
	r := srv.r.Group("/grades")
	r.Use(srv.Authentication())
	{
		r.GET("", srv.HandleGetAllGrades())
	}
}

func (srv *Server) HandleGetAllGrades() gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []*repository.User
		if _, err := srv.dc.GetAll(c.Request.Context(), datastore.NewQuery(repository.KindUser), &users); err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		gpaJsons := make([]*models.Gpa, len(users))
		for i, user := range users {
			gpaJsons[i] = &models.Gpa{Gpa: user.Gpa}
		}
		c.JSON(http.StatusOK, gpaJsons)
	}
}

package server

import (
	"lab-assignment-system-backend/repository"
	"lab-assignment-system-backend/server/models"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

func (srv *Server) LabsRouter() {
	gradesRouter := srv.r.Group("/labs")
	gradesRouter.Use(srv.Authentication())
	{
		gradesRouter.GET("", srv.HandleGetLabs())
	}
}

func (srv *Server) HandleGetLabs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		repoLabs := make([]*repository.Lab, 0, 100)
		if _, err := srv.dc.GetAll(ctx, datastore.NewQuery(repository.KindLab), repoLabs); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		labs := make([]*models.Lab, len(repoLabs))
		for i, repoLab := range repoLabs {
			labs[i] = &models.Lab{
				ID:       repoLab.ID,
				Name:     repoLab.Name,
				Capacity: repoLab.Capacity,
			}
		}

		c.JSON(http.StatusOK, &models.LabList{Labs: labs})
	}
}

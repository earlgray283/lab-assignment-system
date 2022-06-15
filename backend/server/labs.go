package server

import (
	"lab-assignment-system-backend/repository"
	"lab-assignment-system-backend/server/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (srv *Server) LabsRouter() {
	gradesRouter := srv.r.Group("/labs")
	gradesRouter.Use(srv.Authentication())
	{
		gradesRouter.GET("", srv.HandleGetAllLabs())
	}
}

func (srv *Server) HandleGetAllLabs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var repoLabs []*repository.Lab
		labIds := strings.Split(c.Query("labIds"), "+")
		repoLabs, ok, err := repository.FetchAllLabs(ctx, srv.dc, labIds)
		if err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !ok {
			AbortWithErrorJSON(c, NewError(http.StatusNotFound, "no such lab"))
			return
		}

		labs := make([]*models.Lab, len(repoLabs))
		for i, repoLab := range repoLabs {
			labs[i] = &models.Lab{
				ID:           repoLab.ID,
				Name:         repoLab.Name,
				Capacity:     repoLab.Capacity,
				FirstChoice:  repoLab.FirstChoice,
				SecondChoice: repoLab.SecondChoice,
				ThirdChoice:  repoLab.ThirdChice,
			}
		}

		c.JSON(http.StatusOK, &models.LabList{Labs: labs})
	}
}

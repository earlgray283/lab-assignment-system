package server

import (
	"lab-assignment-system-backend/repository"
	"lab-assignment-system-backend/server/models"
	"net/http"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

const LabsCapacity = 1024

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
		if len(labIds) == 0 {
			repoLabs2 := make([]*repository.Lab, 0, LabsCapacity)
			if _, err := srv.dc.GetAll(ctx, datastore.NewQuery(repository.KindLab), &repoLabs2); err != nil {
				srv.logger.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			repoLabs = repoLabs2
		} else {
			repoLabs2 := make([]*repository.Lab, len(labIds))
			repoKeys := make([]*datastore.Key, len(labIds))
			for i, labId := range labIds {
				repoKeys[i] = repository.NewLabKey(labId)
			}
			if err := srv.dc.GetMulti(ctx, repoKeys, repoLabs2); err != nil {
				srv.logger.Println(err)
				if merr, ok := err.(datastore.MultiError); ok {
					for _, err := range merr {
						if err == datastore.ErrNoSuchEntity {
							c.AbortWithStatus(http.StatusNotFound)
							return
						}
					}
				}
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			repoLabs = repoLabs2
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

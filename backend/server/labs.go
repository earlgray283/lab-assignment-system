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
		labIds, _ := c.GetQueryArray("labId")
		keys := make([]*datastore.Key, 0)
		if len(labIds) == 0 {
			keys2, err := srv.dc.GetAll(ctx, datastore.NewQuery(repository.KindLab).KeysOnly(), nil)
			if err != nil {
				srv.logger.Println("check", err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			keys = keys2
		} else {
			for _, labId := range labIds {
				keys = append(keys, repository.NewLabKey(labId))
			}
		}

		repoLabs := make([]repository.Lab, len(keys))
		if err := srv.dc.GetMulti(ctx, keys, repoLabs); err != nil {
			srv.logger.Println(err)
			if merr, ok := err.(datastore.MultiError); !ok {
				for _, err := range merr {
					if err == datastore.ErrNoSuchEntity {
						AbortWithErrorJSON(c, NewError(http.StatusNotFound, "no such lab"))
						return
					}
				}
			}
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		labs := make([]models.Lab, len(keys))
		baseQuery := datastore.NewQuery(repository.KindUser)
		for i, repoLab := range repoLabs {
			var err error
			labs[i] = models.Lab{
				ID:       repoLab.ID,
				Name:     repoLab.Name,
				Capacity: repoLab.Capacity,
			}
			labs[i].FirstChoice, err = srv.dc.Count(ctx, baseQuery.Filter("Lab1 = ", repoLab.ID))
			if err != nil {
				srv.logger.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			labs[i].SecondChoice, err = srv.dc.Count(ctx, baseQuery.Filter("Lab2 = ", repoLab.ID))
			if err != nil {
				srv.logger.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			labs[i].ThirdChoice, err = srv.dc.Count(ctx, baseQuery.Filter("Lab3 = ", repoLab.ID))
			if err != nil {
				srv.logger.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		c.JSON(http.StatusOK, labs)
	}
}

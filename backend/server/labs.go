package server

import (
	"fmt"
	"lab-assignment-system-backend/repository"
	"lab-assignment-system-backend/server/models"
	"net/http"
	"strings"

	"cloud.google.com/go/datastore"
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

		labIds := strings.Split(c.Query("labIds"), "+")
		var optFields []string
		if optFieldsText, ok := c.GetQuery("optFields"); ok {
			optFields = strings.Split(optFieldsText, "+")
		}
		// TODO: definition type を使うなりする
		for _, optField := range optFields {
			if optField != "grade" {
				AbortWithErrorJSON(c, NewError(http.StatusBadRequest, fmt.Sprintln("optField", optField, "is not supported")))
				return
			}
		}

		var repoLabs []*repository.Lab
		if len(labIds) == 0 {
			if _, err := srv.dc.GetAll(ctx, datastore.NewQuery(repository.KindLab), &repoLabs); err != nil {
				srv.logger.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else {
			repoLabs2, ok, err := repository.FetchAllLabs(ctx, srv.dc, labIds)
			if err != nil {
				srv.logger.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			if !ok {
				AbortWithErrorJSON(c, NewError(http.StatusNotFound, "no such lab"))
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
			for _, optField := range optFields {
				switch optField {
				case "grade":
					labs[i].Grades = srv.gpaWorker.Get(repoLab.ID)
				}
			}
		}

		c.JSON(http.StatusOK, &models.LabList{Labs: labs})
	}
}

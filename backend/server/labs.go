package server

import (
	"fmt"
	"lab-assignment-system-backend/lib"
	"lab-assignment-system-backend/repository"
	"lab-assignment-system-backend/server/models"
	"net/http"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

func (srv *Server) LabsRouter() {
	r := srv.r.Group("/labs")
	{
		r.GET("", srv.HandleGetAllLabs())
	}
}

func (srv *Server) HandleGetAllLabs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		labIdsText, ok := c.GetQuery("labIds")
		var labIds []string
		if ok {
			labIds = strings.Split(labIdsText, "+")
		}
		var optFields []string
		if optFieldsText, ok := c.GetQuery("optFields"); ok {
			optFields = strings.Split(optFieldsText, "+")
		}
		// TODO: definition type を使うなりする
		for _, optField := range optFields {
			if optField != "grade" {
				lib.AbortWithErrorJSON(c, lib.NewError(http.StatusBadRequest, fmt.Sprintln("optField", optField, "is not supported")))
				return
			}
		}

		var repoLabs []*repository.Lab
		if len(labIds) == 0 {
			if _, err := srv.dc.GetAll(ctx, datastore.NewQuery(repository.KindLab), &repoLabs); err != nil {
				srv.logger.Printf("%+v\n", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else {
			repoLabs2, ok, err := repository.FetchAllLabs(ctx, srv.dc, labIds)
			if err != nil {
				srv.logger.Printf("%+v\n", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			if !ok {
				lib.AbortWithErrorJSON(c, lib.NewError(http.StatusNotFound, "no such lab"))
				return
			}
			repoLabs = repoLabs2
		}

		labs := make([]*models.Lab, len(repoLabs))
		for i, repoLab := range repoLabs {
			labGpa := srv.labsChecker.GetLabGpa(repoLab.ID)
			labs[i] = &models.Lab{
				ID:           repoLab.ID,
				Name:         repoLab.Name,
				Capacity:     repoLab.Capacity,
				FirstChoice:  len(labGpa.Gpas1),
				SecondChoice: len(labGpa.Gpas2),
				ThirdChoice:  len(labGpa.Gpas3),
			}
			for _, optField := range optFields {
				switch optField {
				case "grade":
					labs[i].Grades = labGpa
				}
			}
		}

		c.JSON(http.StatusOK, &models.LabList{Labs: labs})
	}
}

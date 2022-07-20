package server

import (
	"encoding/csv"
	"fmt"
	"lab-assignment-system-backend/lib"
	"lab-assignment-system-backend/repository"
	"lab-assignment-system-backend/server/models"
	"net/http"
	"sort"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

func (srv *Server) LabsRouter() {
	r := srv.r.Group("/labs")
	{
		r.GET("", srv.HandleGetAllLabs())
	}
	srv.r.POST("/labs/confirm", srv.HandlePostConfirmLabs())
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
				ID:              repoLab.ID,
				Name:            repoLab.Name,
				Capacity:        repoLab.Capacity,
				ConfirmedNumber: repoLab.ConfirmedNumber,
				FirstChoice:     len(labGpa.Gpas1),
				SecondChoice:    len(labGpa.Gpas2),
				ThirdChoice:     len(labGpa.Gpas3),
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

func (srv *Server) HandlePostConfirmLabs() gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []*repository.User
		if _, err := srv.dc.GetAll(c.Request.Context(), datastore.NewQuery(repository.KindUser), &users); err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var labs []*repository.Lab
		if _, err := srv.dc.GetAll(c.Request.Context(), datastore.NewQuery(repository.KindLab), &labs); err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		sort.Slice(users, func(i, j int) bool { return users[i].Gpa > users[j].Gpa })

		userKeys := make([]*datastore.Key, 0, len(users))
		labMap := lib.NewMapFromSlice(lib.MapSlice(labs, func(lab *repository.Lab) string { return lab.ID }), labs)
		for _, user := range users {
			if user.Lab1 == nil {
				continue
			}
			if labMap[*user.Lab1].Capacity == labMap[*user.Lab1].ConfirmedNumber {
				continue
			}
			user.ConfirmedLab = user.Lab1
			labMap[*user.Lab1].ConfirmedNumber++
			userKeys = append(userKeys, repository.NewUserKey(user.UID))
		}
		newLabKeys := make([]*datastore.Key, 0, len(labs))
		newLabs := make([]*repository.Lab, 0, len(labs))
		for _, lab := range labMap {
			newLabKeys = append(newLabKeys, repository.NewLabKey(lab.ID))
			newLabs = append(newLabs, lab)
		}

		if _, err := srv.dc.RunInTransaction(c.Request.Context(), func(tx *datastore.Transaction) error {
			if _, err := tx.PutMulti(userKeys, users); err != nil {
				return err
			}
			if _, err := tx.PutMulti(newLabKeys, newLabs); err != nil {
				return err
			}
			return nil
		}); err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Writer.Header().Set("Content-Type", "text/csv")
		csvw := csv.NewWriter(c.Writer)
		for _, user := range users {
			confirmedLab := "undefined"
			if user.ConfirmedLab != nil {
				confirmedLab = *user.ConfirmedLab
			}
			csvw.Write([]string{user.UID, confirmedLab})
		}
		csvw.Flush()
	}
}

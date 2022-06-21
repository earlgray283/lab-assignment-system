package server

import (
	"lab-assignment-system-backend/lib"
	"lab-assignment-system-backend/repository"
	"lab-assignment-system-backend/server/models"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func (srv *Server) GradesRouter() {
	r := srv.r.Group("/grades")
	r.Use(srv.Authentication())
	{
		r.GET("", srv.HandleGetAllGrades())
		r.POST("/generate-token", srv.HandleGenerateToken())
	}

	// authentication middleware を適用しない
	srv.r.POST("/grades", srv.HandlePostGrade())
}

func (srv *Server) HandleGetAllGrades() gin.HandlerFunc {
	return func(c *gin.Context) {
		gpas := srv.gpaWorker.GetGpas()
		gpaJsons := make([]*models.Gpa, len(gpas))
		for i, gpa := range gpas {
			gpaJsons[i] = &models.Gpa{Gpa: gpa}
		}
		c.JSON(http.StatusOK, gpaJsons)
	}
}

func (srv *Server) HandleGenerateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authToken := c.MustGet("authToken").(*auth.Token)

		token := lib.MakeRandomString(32)
		now := time.Now()
		gradeRequestToken := &repository.RegisterToken{
			UID:       authToken.UID,
			Token:     token,
			Expires:   now.Add(time.Hour),
			CreatedAt: now,
		}
		if _, err := srv.dc.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
			if _, err := tx.Put(repository.NewRegisterTokenKey(token), gradeRequestToken); err != nil {
				return err
			}
			return nil
		}); err != nil {
			srv.logger.Printf("%+v\n", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gradeRequestToken)
	}
}

func (srv *Server) HandlePostGrade() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		token := c.Request.Header.Get("register-token")
		if token == "" {
			lib.AbortWithErrorJSON(c, lib.NewError(http.StatusBadRequest, "You must generate register-token"))
			return
		}
		var grade models.Grade
		if err := c.BindJSON(&grade); err != nil {
			srv.logger.Printf("%+v\n", err)
			return
		}

		var registerToken repository.RegisterToken
		if err := srv.dc.Get(ctx, repository.NewRegisterTokenKey(token), &registerToken); err != nil {
			srv.logger.Printf("%+v\n", err)
			if err == datastore.ErrNoSuchEntity {
				lib.AbortWithErrorJSON(c, lib.NewError(http.StatusBadRequest, "the register-token was not found"))
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}
		if time.Now().After(registerToken.Expires) {
			lib.AbortWithErrorJSON(c, lib.NewError(http.StatusBadRequest, "invalid register-token"))
			return
		}

		var user repository.User
		userKey := repository.NewUserKey(registerToken.UID)
		if err := srv.dc.Get(c, userKey, &user); err != nil {
			srv.logger.Printf("%+v\n", err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		if user.Gpa == nil {
			gpa := lib.CalculateGpa(grade.Grades, &lib.CalculateGpaOption{
				Until:             time.Date(time.Now().Year(), 3, 31, 23, 59, 59, 0, time.Local),
				ExcludeLowerPoint: 60,
			})
			user.Gpa = &gpa
			if _, err := srv.dc.Put(c, userKey, &user); err != nil {
				srv.logger.Printf("%+v\n", err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		} else {
			// 既に成績が存在する場合は処理を行わない
			// 上書きしたい場合は delete リクエストを送ってもらう
			lib.AbortWithErrorJSON(c, lib.NewError(http.StatusConflict, "Grade has already existed. Please delete it and try again."))
		}
	}
}

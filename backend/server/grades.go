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
	gradesRouter := srv.r.Group("/grades")
	gradesRouter.Use(srv.Authentication())
	{
		gradesRouter.GET("/", srv.HandleGetGpas())
		gradesRouter.GET("/me", srv.HandleGetOwnGpa())
		gradesRouter.POST("/generate-token", srv.HandleGenerateToken())
	}

	// authentication middleware を適用しない
	srv.r.POST("/grades", srv.HandlePostGrade())
}

func (srv *Server) HandleGetOwnGpa() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("authToken").(*auth.Token)
		ctx := c.Request.Context()

		var grade repository.Grade
		if err := srv.dc.Get(ctx, repository.NewGradeKey(authToken.UID), &grade); err != nil {
			if err == datastore.ErrNoSuchEntity {
				AbortWithErrorJSON(c, NewError(http.StatusNotFound, "there are no grades"))
			} else {
				srv.logger.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}
		c.JSON(http.StatusOK, &grade)
	}
}
func (srv *Server) HandleGetGpas() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var grades []repository.Grade
		if _, err := srv.dc.GetAll(ctx, datastore.NewQuery(repository.KindGrade), &grades); err != nil {
			if err == datastore.ErrNoSuchEntity {
				AbortWithErrorJSON(c, NewError(http.StatusNotFound, "there are no grades"))
			} else {
				srv.logger.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}
		c.JSON(http.StatusOK, &grades)
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
			srv.logger.Println(err)
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
			AbortWithErrorJSON(c, NewError(http.StatusBadRequest, "You must generate register-token"))
			return
		}
		var grade models.Grade
		if err := c.BindJSON(&grade); err != nil {
			srv.logger.Println(err)
			return
		}

		var registerToken repository.RegisterToken
		if err := srv.dc.Get(ctx, repository.NewRegisterTokenKey(token), &registerToken); err != nil {
			srv.logger.Println(err)
			if err == datastore.ErrNoSuchEntity {
				AbortWithErrorJSON(c, NewError(http.StatusBadRequest, "Invalid register-token"))
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}

		var repoGrade repository.Grade
		gradeKey := repository.NewGradeKey(registerToken.UID)
		if err := srv.dc.Get(c, gradeKey, &repoGrade); err != nil {
			srv.logger.Println(err)
			if err != datastore.ErrNoSuchEntity {
				c.AbortWithStatus(http.StatusInternalServerError)
			} else {
				if _, err := srv.dc.Put(c, gradeKey, &repository.Grade{
					UID:           registerToken.UID,
					StudentName:   grade.StudentName,
					StudentNumber: grade.StudentNumber,
					Gpa: lib.CalculateGpa(grade.Grades, &lib.CalculateGpaOption{
						Until:             time.Date(time.Now().Year(), 3, 31, 23, 59, 59, 0, time.Local),
						ExcludeLowerPoint: 60,
					}),
				}); err != nil {
					srv.logger.Println(err)
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		} else {
			// 既に成績が存在する場合は処理を行わない
			// 上書きしたい場合は delete リクエストを送ってもらう
			srv.logger.Println(err)
			AbortWithErrorJSON(c, NewError(http.StatusConflict, "Grade has already existed. Please delete it and try again."))
		}
	}
}

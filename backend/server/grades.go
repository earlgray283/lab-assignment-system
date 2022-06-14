package server

import (
	"lab-assignment-system-backend/lib"
	"lab-assignment-system-backend/repository"
	"lab-assignment-system-backend/server/models"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

func (srv *Server) GradesRouter() {
	gradesRouter := srv.r.Group("/grades")
	{
		gradesRouter.POST("/", srv.HandlePostGrade())
		gradesRouter.POST("/generate-token", srv.HandleGenerateToken())
		gradesRouter.GET("/gpa", srv.HandleGetGpa())
	}
}

func (srv *Server) HandleGetGpa() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authToken, err := srv.GetAuthToken(c)
		if err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var grade repository.Grade
		if err := srv.dc.Get(ctx, repository.NewGradeKey(authToken.UID), &grade); err != nil {
			srv.logger.Println(err)
			AbortWithErrorJSON(c, NewError(http.StatusNotFound, "there are no grades"))
			return
		}
		c.JSON(http.StatusOK, &grade)
	}
}

func (srv *Server) HandleGenerateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authToken, err := srv.GetAuthToken(c)
		if err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

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

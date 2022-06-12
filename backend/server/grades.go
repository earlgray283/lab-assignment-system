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
		gpa := lib.CalculateGpa(&grade, &lib.CalculateGpaOption{
			Until:             time.Date(time.Now().Year(), time.March, 31, 0, 0, 0, 0, nil),
			ExcludeLowerPoint: ExcludeLowerPoint,
		})
		c.JSON(http.StatusOK, &models.Gpa{Gpa: gpa})
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
		gradeRequestToken := &repository.GradeRequestToken{
			UID:       authToken.UID,
			Token:     token,
			Expires:   now.Add(time.Hour),
			CreatedAt: now,
		}
		if _, err := srv.dc.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
			if _, err := tx.Put(repository.NewGradeRequestTokenKey(token), gradeRequestToken); err != nil {
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
		authToken, err := srv.GetAuthToken(c)
		if err != nil {
			AbortWithErrorJSON(c, NewError(http.StatusUnauthorized, "not logged in"))
		}
		var grade, empty repository.Grade
		if err := c.BindJSON(&grade); err != nil {
			srv.logger.Println(err)
			return
		}
		key := repository.NewGradeKey(authToken.UID)
		err = srv.dc.Get(c, key, &empty)
		if err == nil {
			// 既に成績が存在する場合は処理を行わない
			// 上書きしたい場合は delete リクエストを送ってもらう
			srv.logger.Println("Duplicate post grades request")
			c.AbortWithStatus(http.StatusConflict)
			return
		}
		if err != datastore.ErrNoSuchEntity {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if _, err := srv.dc.Put(c, key, &grade); err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}

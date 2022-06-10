package server

import (
	"lab-assignment-system-backend/lib"
	"lab-assignment-system-backend/repository"
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
	}
}

func (srv *Server) HandleGenerateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		sessionCookie, err := c.Cookie("session")
		if err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		authToken, err := srv.auth.VerifySessionCookie(ctx, sessionCookie)
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

		var grade, empty repository.Grade
		if err := c.BindJSON(&grade); err != nil {
			srv.logger.Println(err)
			return
		}
		key := repository.NewGradeKey(grade.StudentNumber)
		err := srv.dc.Get(c, key, &empty)
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

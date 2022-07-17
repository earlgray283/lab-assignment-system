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

func (srv *Server) UserRouter() {
	r := srv.r.Group("/user")
	r.Use(srv.Authentication())
	{
		r.GET("", srv.HandleGetUser())
		r.PUT("/lab", srv.HandleUpdateLabs())
	}
}

func (srv *Server) HandleUpdateLabs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		user, _ := GetUser(c)
		var userLab models.UserLab
		if err := c.BindJSON(&userLab); err != nil {
			srv.logger.Println(err)
			return
		}

		user.Lab1 = &userLab.Lab1
		user.Lab2 = &userLab.Lab2
		user.Lab3 = &userLab.Lab3
		user.UpdatedAt = lib.PointerOfValue(time.Now())

		userKey := repository.NewUserKey(user.UID)
		if _, err := srv.dc.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
			if _, err := tx.Put(userKey, user); err != nil {
				return err
			}
			return nil
		}); err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if err := srv.labsChecker.SingleRun(); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, &models.User{
			UID:  user.UID,
			Gpa:  user.Gpa,
			Lab1: user.Lab1,
			Lab2: user.Lab2,
			Lab3: user.Lab3,
		})
	}
}

func (srv *Server) HandleGetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		repoUser, err := GetUser(c)
		if err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		user := &models.User{
			UID:  repoUser.UID,
			Gpa:  repoUser.Gpa,
			Lab1: repoUser.Lab1,
			Lab2: repoUser.Lab2,
			Lab3: repoUser.Lab3,
		}
		c.JSON(http.StatusOK, user)
	}
}

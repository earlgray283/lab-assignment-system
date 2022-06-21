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

// TODO: クソ uri をなんとかする
func (srv *Server) UserRouter() {
	r := srv.r.Group("/users")
	r.Use(srv.Authentication())
	{
		r.GET("", srv.HandleGetUser())
		r.PUT("", srv.HandlePutUser())
		r.DELETE("/:uid", srv.HandleDeleteUser()).Use(srv.CheckAdmin())
	}
}

func (srv *Server) HandlePutUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authToken := c.MustGet("authToken").(*auth.Token)
		var newUser models.User
		if err := c.BindJSON(&newUser); err != nil {
			srv.logger.Println(err)
			return
		}
		var user repository.User
		if err := srv.dc.Get(ctx, repository.NewUserKey(authToken.UID), &user); err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		repoNewUser, userKey := repository.NewUser(authToken.UID, user.Email, user.StudentNumber, user.Name, newUser.Lab1, newUser.Lab2, newUser.Lab3, user.Gpa, user.CreatedAt)
		repoNewUser.UpdatedAt = lib.PointerOfValue(time.Now())
		if _, err := srv.dc.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
			if _, err := tx.Put(userKey, repoNewUser); err != nil {
				return err
			}
			return nil
		}); err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}

func (srv *Server) HandleDeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		uid := c.Param("uid")
		userKey := repository.NewUserKey(uid)
		var user repository.User
		if err := srv.dc.Get(ctx, userKey, &user); err != nil {
			srv.logger.Println(err)
			if err == datastore.ErrNoSuchEntity {
				lib.AbortWithErrorJSON(c, lib.NewError(http.StatusBadRequest, "no such user"))
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}
		if _, err := srv.dc.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
			if err := tx.Delete(userKey); err != nil {
				return err
			}
			if err := srv.auth.DeleteUser(ctx, uid); err != nil {
				return err
			}
			return nil
		}); err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}

func (srv *Server) HandleGetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authToken := c.MustGet("authToken").(*auth.Token)
		if authToken == nil {
			srv.logger.Println("context value authToken was nil")
			lib.AbortWithErrorJSON(c, lib.NewError(http.StatusUnauthorized, "not logged in"))
			return
		}
		var user repository.User
		if err := srv.dc.Get(ctx, repository.NewUserKey(authToken.UID), &user); err != nil {
			srv.logger.Printf("%+v\n", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, &models.User{
			UID:           user.UID,
			Email:         user.Email,
			StudentNumber: user.StudentNumber,
			Name:          user.Name,
			Gpa:           user.Gpa,
			Lab1:          user.Lab1,
			Lab2:          user.Lab2,
			Lab3:          user.Lab3,
		})
	}
}

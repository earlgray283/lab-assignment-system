package server

import (
	"lab-assignment-system-backend/lib"
	"lab-assignment-system-backend/repository"
	"lab-assignment-system-backend/server/models"
	"net/http"

	"cloud.google.com/go/datastore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func (srv *Server) UserRouter() {
	gradesRouter := srv.r.Group("/users")
	gradesRouter.Use(srv.Authentication())
	{
		gradesRouter.GET("", srv.HandleGetUser())
		gradesRouter.DELETE("/:uid", srv.HandleDeleteUser())
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
		keys := []*datastore.Key{
			repository.NewLabKey(user.Lab1),
			repository.NewLabKey(user.Lab2),
			repository.NewLabKey(user.Lab3),
		}
		labs := make([]*repository.Lab, 3)
		if err := srv.dc.GetMulti(ctx, keys, labs); err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if labs[0].FirstChoice > 0 {
			labs[0].FirstChoice--
		}
		if labs[1].SecondChoice > 0 {
			labs[1].SecondChoice--
		}
		if labs[2].ThirdChice > 0 {
			labs[2].ThirdChice--
		}
		if _, err := srv.dc.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
			mutations := make([]*datastore.Mutation, 0, 4)
			for i := 0; i < 3; i++ {
				mutations = append(mutations, datastore.NewUpdate(keys[i], labs[i]))
			}
			mutations = append(mutations, datastore.NewDelete(userKey))
			if _, err := srv.dc.Mutate(ctx, mutations...); err != nil {
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

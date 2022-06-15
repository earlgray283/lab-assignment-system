package server

import (
	"lab-assignment-system-backend/repository"
	"lab-assignment-system-backend/server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (srv *Server) UserRouter() {
	gradesRouter := srv.r.Group("/users")
	{
		gradesRouter.GET("/:uid", srv.HandleGetUser())
	}
}

func (srv *Server) HandleGetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authToken, err := srv.GetAuthToken(c)
		if err != nil {
			srv.logger.Println(err)
			AbortWithErrorJSON(c, NewError(http.StatusUnauthorized, "not logged in"))
			return
		}
		uid := c.Param("uid")
		if authToken.UID != uid {
			// TODO: admin 権限で見れるようにする
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		var user repository.User
		if err := srv.dc.Get(ctx, repository.NewUserKey(uid), &user); err != nil {
			srv.logger.Println(err)
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

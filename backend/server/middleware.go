package server

import (
	"errors"
	"lab-assignment-system-backend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (srv *Server) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionValue, err := c.Cookie("session")
		if err != nil {
			srv.logger.Printf("%+v\n", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var session repository.Session
		if err := srv.dc.Get(c.Request.Context(), repository.NewSessionKey(sessionValue), &session); err != nil {
			srv.logger.Printf("%+v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		var user repository.User
		if err := srv.dc.Get(c.Request.Context(), repository.NewUserKey(session.UID), &user); err != nil {
			srv.logger.Printf("%+v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Set("user", &user)
		c.Next()
	}
}

func GetUser(c *gin.Context) (*repository.User, error) {
	v, exists := c.Get("user")
	if !exists {
		return nil, errors.New("context value user was not found")
	}
	user, ok := v.(*repository.User)
	if !ok {
		return nil, errors.New("cannot convert any to *repository.User")
	}
	return user, nil
}

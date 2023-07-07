package middleware

import (
	"errors"
	"lab-assignment-system-backend/server/domain/entity"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

func Authentication(dsClient *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionValue, err := c.Cookie("session")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var session entity.Session
		if err := dsClient.Get(c.Request.Context(), entity.NewSessionKey(sessionValue), &session); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		var user entity.User
		if err := dsClient.Get(c.Request.Context(), entity.NewUserKey(session.UID), &user); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Set("user", &user)
		c.Next()
	}
}

func GetUser(c *gin.Context) (*entity.User, error) {
	v, exists := c.Get("user")
	if !exists {
		return nil, errors.New("context value user was not found")
	}
	user, ok := v.(*entity.User)
	if !ok {
		return nil, errors.New("cannot convert any to *entity.User")
	}
	return user, nil
}

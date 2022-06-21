package server

import (
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func (srv *Server) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken, err := srv.GetAuthToken(c)
		if err != nil {
			srv.logger.Printf("%+v\n", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("authToken", authToken)
		c.Next()
	}
}

func (srv *Server) CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("authToken").(*auth.Token)
		user, err := srv.auth.GetUser(c.Request.Context(), authToken.UID)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims := user.CustomClaims
		if claimValue, ok := claims["admin"]; ok {
			isAdmin := claimValue.(bool)
			if !isAdmin {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		} else {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

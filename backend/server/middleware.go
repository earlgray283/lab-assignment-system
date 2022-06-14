package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (srv *Server) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken, err := srv.GetAuthToken(c)
		if err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("authToken", authToken)
		c.Next()
	}
}

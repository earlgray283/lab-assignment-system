package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors(allowOrigins []string) gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = append([]string{"http://localhost:3000"}, allowOrigins...)
	config.AllowCredentials = true
	return cors.New(config)
}

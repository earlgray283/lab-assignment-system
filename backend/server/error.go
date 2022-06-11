package server

import "github.com/gin-gonic/gin"

type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewError(code int, message string) *Error {
	return &Error{code, message}
}

func AbortWithErrorJSON(c *gin.Context, errJson *Error) {
	c.AbortWithStatusJSON(errJson.Code, errJson)
}

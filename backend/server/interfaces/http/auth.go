package http

import (
	"lab-assignment-system-backend/server/domain/models"
	"lab-assignment-system-backend/server/lib"
	"lab-assignment-system-backend/server/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	interactor *usecases.AuthInteractor
}

func NewAuthController(interactor *usecases.AuthInteractor) *AuthController {
	return &AuthController{interactor}
}

func (c *AuthController) Login(gc *gin.Context) {
	var payload models.SigninPayload
	if err := gc.ShouldBindJSON(&payload); err != nil {
		gc.AbortWithStatusJSON(http.StatusBadRequest, lib.NewBadRequestError("invalid payload"))
	}

	resp, cookie, err := c.interactor.Login(gc.Request.Context(), payload.UID)
	if err != nil {
		err := err.(*lib.Error)
		gc.AbortWithStatusJSON(err.Code, err)
		return
	}
	http.SetCookie(gc.Writer, cookie)

	gc.JSON(http.StatusOK, resp)
}

func (c *AuthController) Logout(gc *gin.Context) {
	cookie, err := gc.Request.Cookie("session")
	if err != nil {
		gc.AbortWithStatusJSON(http.StatusBadRequest, lib.NewBadRequestError("no such session cookie"))
		return
	}
	c.interactor.Logout(cookie)
	http.SetCookie(gc.Writer, cookie)
}

package http

import (
	"lab-assignment-system-backend/server/api/middleware"
	"lab-assignment-system-backend/server/domain/models"
	"lab-assignment-system-backend/server/lib"
	"lab-assignment-system-backend/server/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	interactor *usecases.UsersInteractor
}

func NewUsersController(interactor *usecases.UsersInteractor) *UsersController {
	return &UsersController{interactor}
}

func (c *UsersController) UpdateUser(gc *gin.Context) {
	user, _ := middleware.GetUser(gc)
	if user.UID == "test" && !lib.IsDebug() {
		gc.AbortWithStatusJSON(http.StatusForbidden, "テストユーザーは編集できません")
		return
	}

	var payload models.UpdateUserPayload
	if err := gc.ShouldBindJSON(&payload); err != nil {
		err := err.(*lib.Error)
		gc.AbortWithStatusJSON(err.Code, err)
		return
	}

	resp, err := c.interactor.UpdateUser(gc.Request.Context(), user, &payload)
	if err != nil {
		err := err.(*lib.Error)
		gc.AbortWithStatusJSON(err.Code, err)
		return
	}

	gc.JSON(200, resp)
}

func (c *UsersController) GetUserMe(gc *gin.Context) {
	user, _ := middleware.GetUser(gc)

	gc.JSON(200, &models.User{
		UID:          user.UID,
		Gpa:          user.Gpa,
		Role:         string(user.Role),
		WishLab:      user.WishLab,
		ConfirmedLab: user.ConfirmedLab,
		Year:         user.Year,
		Reason:       user.Reason,
	})
}

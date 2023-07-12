package http

import (
	"lab-assignment-system-backend/server/api/middleware"
	"lab-assignment-system-backend/server/domain/entity"
	"lab-assignment-system-backend/server/domain/models"
	"lab-assignment-system-backend/server/usecases"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	interactor *usecases.AdminInteractor
}

func NewAdminController(interactor *usecases.AdminInteractor) *AdminController {
	return &AdminController{interactor}
}

func (c *AdminController) FinalDecision(gc *gin.Context) {
	user, _ := middleware.GetUser(gc)
	if user.Role != entity.RoleAdmin {
		gc.AbortWithStatusJSON(403, "権限がありません")
		return
	}

	var payload models.FinalDecisionPayload
	if err := gc.ShouldBindJSON(&payload); err != nil {
		gc.AbortWithStatusJSON(400, "invalid payload")
		return
	}

	resp, err := c.interactor.FinalDecision(gc.Request.Context(), payload.Year)
	if err != nil {
		gc.AbortWithStatusJSON(500, err)
		return
	}
	gc.JSON(200, resp)
}

package http

import (
	"io"
	"lab-assignment-system-backend/server/api/middleware"
	"lab-assignment-system-backend/server/domain/entity"
	"lab-assignment-system-backend/server/domain/models"
	"lab-assignment-system-backend/server/lib"
	"lab-assignment-system-backend/server/usecases"
	"strconv"

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

func (c *AdminController) GetCSV(gc *gin.Context) {
	user, _ := middleware.GetUser(gc)
	if user.Role != entity.RoleAdmin {
		gc.AbortWithStatusJSON(403, "権限がありません")
		return
	}

	year, err := strconv.ParseInt(gc.Query("year"), 10, 64)
	if err != nil {
		gc.AbortWithStatusJSON(400, "year must be set")
		return
	}

	r, err := c.interactor.GetCSV(gc.Request.Context(), int(year))
	if err != nil {
		err := err.(*lib.Error)
		gc.AbortWithStatusJSON(err.Code, err)
		return
	}

	gc.Writer.Header().Set("Content-Type", "text/csv")
	_, _ = io.Copy(gc.Writer, r)
}

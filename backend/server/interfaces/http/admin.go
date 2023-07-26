package http

import (
	"archive/zip"
	"io"
	"lab-assignment-system-backend/server/api/middleware"
	"lab-assignment-system-backend/server/domain/entity"
	"lab-assignment-system-backend/server/domain/models"
	"lab-assignment-system-backend/server/lib"
	"lab-assignment-system-backend/server/usecases"
	"net/http"
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
		gc.AbortWithStatusJSON(403, lib.NewBadRequestError("権限がありません"))
		return
	}

	var payload models.FinalDecisionPayload
	if err := gc.ShouldBindJSON(&payload); err != nil {
		gc.AbortWithStatusJSON(400, lib.NewBadRequestError("invalid payload"))
		return
	}

	labCSV, userCSV, err := c.interactor.FinalDecision(gc.Request.Context(), payload.Year)
	if err != nil {
		gc.AbortWithStatusJSON(500, err)
		return
	}

	zw := zip.NewWriter(gc.Writer)
	labw, err := zw.Create("lab.csv")
	if err != nil {
		gc.AbortWithStatusJSON(500, err)
		return
	}
	_, err = io.Copy(labw, labCSV)
	if err != nil {
		gc.AbortWithStatusJSON(500, err)
		return
	}
	userw, err := zw.Create("user.csv")
	if err != nil {
		gc.AbortWithStatusJSON(500, err)
		return
	}
	_, err = io.Copy(userw, userCSV)
	if err != nil {
		gc.AbortWithStatusJSON(500, err)
		return
	}
	zw.Close()
	gc.Writer.Header().Set("Content-Type", "application/zip")
	gc.Writer.WriteHeader(http.StatusOK)
}

func (c *AdminController) FinalDecisionDryRun(gc *gin.Context) {
	user, _ := middleware.GetUser(gc)
	if user.Role != entity.RoleAdmin {
		gc.AbortWithStatusJSON(403, lib.NewBadRequestError("権限がありません"))
		return
	}

	var payload models.FinalDecisionPayload
	if err := gc.ShouldBindJSON(&payload); err != nil {
		gc.AbortWithStatusJSON(400, lib.NewBadRequestError("invalid payload"))
		return
	}

	labCSV, userCSV, err := c.interactor.FinalDecisionDryRun(gc.Request.Context(), payload.Year)
	if err != nil {
		gc.AbortWithStatusJSON(500, err)
		return
	}

	zw := zip.NewWriter(gc.Writer)
	labw, err := zw.Create("csv/lab.csv")
	if err != nil {
		gc.AbortWithStatusJSON(500, err)
		return
	}
	_, err = io.Copy(labw, labCSV)
	if err != nil {
		gc.AbortWithStatusJSON(500, err)
		return
	}
	userw, err := zw.Create("csv/user.csv")
	if err != nil {
		gc.AbortWithStatusJSON(500, err)
		return
	}
	_, err = io.Copy(userw, userCSV)
	if err != nil {
		gc.AbortWithStatusJSON(500, err)
		return
	}
	zw.Close()
	gc.Writer.Header().Set("Content-Type", "application/zip")
	gc.Writer.Header().Set("Content-Disposition", "attachment; filename='csv.zip'")
	gc.Writer.WriteHeader(http.StatusOK)
}

func (c *AdminController) GetCSV(gc *gin.Context) {
	user, _ := middleware.GetUser(gc)
	if user.Role != entity.RoleAdmin {
		gc.AbortWithStatusJSON(http.StatusForbidden, lib.NewBadRequestError("権限がありません"))
		return
	}

	year, err := strconv.ParseInt(gc.Query("year"), 10, 64)
	if err != nil {
		gc.AbortWithStatusJSON(http.StatusBadRequest, lib.NewBadRequestError("year must be set"))
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

func (c *AdminController) CreateUsers(gc *gin.Context) {
	user, _ := middleware.GetUser(gc)
	if user.Role != entity.RoleAdmin {
		gc.AbortWithStatusJSON(http.StatusForbidden, lib.NewBadRequestError("権限がありません"))
		return
	}

	var payload models.CreateUsersPayload
	if err := gc.ShouldBindJSON(&payload); err != nil {
		gc.AbortWithStatusJSON(http.StatusBadRequest, lib.NewBadRequestError(err.Error()))
		return
	}

	resp, err := c.interactor.CreateUsers(gc.Request.Context(), &payload)
	if err != nil {
		err := err.(*lib.Error)
		gc.AbortWithStatusJSON(err.Code, err)
		return
	}
	gc.JSON(http.StatusOK, resp)
}

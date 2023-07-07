package http

import (
	"lab-assignment-system-backend/server/lib"
	"lab-assignment-system-backend/server/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GradesController struct {
	interactor *usecases.GradesInteractor
}

func NewGradesController(interactor *usecases.GradesInteractor) *GradesController {
	return &GradesController{interactor}
}

func (c *GradesController) ListGrades(gc *gin.Context) {
	var year int
	if yearText, ok := gc.GetQuery("year"); ok {
		year2, err := strconv.Atoi(yearText)
		if err != nil {
			lib.NewBadRequestError("year must be a number")
			return
		}
		year = year2
	}
	resp, err := c.interactor.ListGrades(gc.Request.Context(), year)
	if err != nil {
		err := err.(*lib.Error)
		gc.AbortWithStatusJSON(err.Code, err)
		return
	}
	gc.JSON(http.StatusOK, resp)
}

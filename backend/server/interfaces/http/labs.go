package http

import (
	"lab-assignment-system-backend/server/lib"
	"lab-assignment-system-backend/server/usecases"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type LabsController struct {
	interactor *usecases.LabsInteractor
}

func NewLabsController(interactor *usecases.LabsInteractor) *LabsController {
	return &LabsController{interactor}
}

func (c *LabsController) ListLabs(gc *gin.Context) {
	var year int
	if yearText, ok := gc.GetQuery("year"); ok {
		year2, err := strconv.Atoi(yearText)
		if err != nil {
			lib.NewBadRequestError("year must be a number")
			return
		}
		year = year2
	} else {
		lib.NewBadRequestError("year is required")
		return
	}
	opts := make([]usecases.ListLabsOptionFunc, 0)
	if labIdsText, ok := gc.GetQuery("labIds"); ok {
		opts = append(opts, usecases.WithLabIDs(strings.Split(labIdsText, "+")))
	}
	if optFieldsText, ok := gc.GetQuery("optFields"); ok {
		opts = append(opts, usecases.WithOptFieldSet(strings.Split(optFieldsText, "+")))
	}

	resp, err := c.interactor.ListLabs(gc.Request.Context(), year, opts...)
	if err != nil {
		err := err.(*lib.Error)
		gc.AbortWithStatusJSON(err.Code, err)
		return
	}

	gc.JSON(http.StatusOK, resp)
}

// TODO: CSV のダウンロード機能

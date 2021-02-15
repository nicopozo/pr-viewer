package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	viewercontext "github.com/nicopozo/pr-viewer/internal/context"
	"github.com/nicopozo/pr-viewer/internal/model"
	"github.com/nicopozo/pr-viewer/internal/service"
)

type PRController interface {
	GetPRs(ginCtx *gin.Context)
}

type prController struct {
	service service.PRService
}

func NewPRController(service service.PRService) (PRController, error) {
	if service == nil {
		return nil, fmt.Errorf("service can not be nil")
	}

	return &prController{service: service}, nil
}

func (controller *prController) GetPRs(ginCtx *gin.Context) {
	reqContext := viewercontext.New(ginCtx)
	logger := viewercontext.Logger(reqContext)

	userType := ginCtx.Request.URL.Query().Get("user_type")
	if userType != "owner" && userType != "reviewer" {
		logger.Error(controller, nil,
			fmt.Errorf("user_type parameter required, only 'owner' or 'reviewer' are valid values"),
			"error getting pull request")
		errorResult := model.NewError(model.ValidationError,
			"user_type parameter required, only 'owner' or 'reviewer' are valid values")

		ginCtx.JSON(http.StatusBadRequest, errorResult)

		return
	}
}

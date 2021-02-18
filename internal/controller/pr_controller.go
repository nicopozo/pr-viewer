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
	GetUsername(ginCtx *gin.Context)
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

func (controller *prController) GetUsername(ginCtx *gin.Context) {
	reqContext := viewercontext.New(ginCtx)
	logger := viewercontext.Logger(reqContext)

	token := ginCtx.Request.URL.Query().Get("token")
	if token == "" {
		logger.Error(controller, nil,
			fmt.Errorf("token parameter required"), "error getting pull request")
		errorResult := model.NewError(model.ValidationError,
			"token parameter required")

		ginCtx.JSON(http.StatusBadRequest, errorResult)

		return
	}

	result, err := controller.service.GetUser(reqContext, token)
	if err != nil {
		logger.Error(controller, nil, err, "Failed to get user")

		errorResult := model.NewError(model.InternalError, "Error occurred getting user. %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, errorResult)

		return
	}

	ginCtx.JSON(http.StatusOK, result)
}

func (controller *prController) GetPRs(ginCtx *gin.Context) {
	reqContext := viewercontext.New(ginCtx)
	logger := viewercontext.Logger(reqContext)

	userType := ginCtx.Request.URL.Query().Get("user_type")
	if userType != "" && userType != "owner" && userType != "reviewer" {
		logger.Error(controller, nil,
			fmt.Errorf("only blank, 'owner' or 'reviewer' are valid values"),
			"error getting pull request")
		errorResult := model.NewError(model.ValidationError, "only blank, 'owner' or 'reviewer' are valid values")

		ginCtx.JSON(http.StatusBadRequest, errorResult)

		return
	}

	token := ginCtx.Request.URL.Query().Get("token")
	if token == "" {
		logger.Error(controller, nil,
			fmt.Errorf("token parameter required"), "error getting pull request")
		errorResult := model.NewError(model.ValidationError,
			"token parameter required")

		ginCtx.JSON(http.StatusBadRequest, errorResult)

		return
	}

	result, err := controller.service.GetPRs(reqContext, userType, token)
	if err != nil {
		logger.Error(controller, nil, err, "Failed to search rules")

		errorResult := model.NewError(model.InternalError, "Error occurred when searching rules. %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, errorResult)

		return
	}

	ginCtx.JSON(http.StatusOK, result)
}

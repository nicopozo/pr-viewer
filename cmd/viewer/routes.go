package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicopozo/pr-viewer/internal/controller"
	"github.com/nicopozo/pr-viewer/internal/github"
	"github.com/nicopozo/pr-viewer/internal/service"
	"github.com/nicopozo/pr-viewer/internal/utils/clients"
)

func mapRoutes(router *gin.Engine) {
	router.Static("/reconciliations/pr-viewer/admin", "/Users/npozo/Proyectos/pr-viewer/web/dist")

	ctrl := newPRController()
	router.GET("/pr-viewer/pull-requests", ctrl.GetPRs)
	router.GET("/pr-viewer/username", ctrl.GetUsername)

	router.GET("/ping", ping)
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func newPRController() controller.PRController {
	httpClient := clients.NewHTTPClient(clients.NewHTTPSettings())

	githubClient, err := github.NewClient(&httpClient)
	if err != nil {
		panic(err)
	}

	svc, err := service.NewGithubPRService(githubClient)
	if err != nil {
		panic(err)
	}

	ctrl, err := controller.NewPRController(svc)
	if err != nil {
		panic(err)
	}

	return ctrl
}

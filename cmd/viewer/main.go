package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nicopozo/pr-viewer/internal/model"
)

func main() {
	router := gin.New()
	router.NoRoute(noRouteHandler)

	router.Use(cors.Default())

	mapRoutes(router)

	if err := router.Run(":8081"); err != nil {
		panic(err.Error())
	}
}

func noRouteHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, model.NewError(model.ResourceNotFoundError, "no handler found for path"))
}

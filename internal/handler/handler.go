package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koma2211/ayan-capital_task/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			h.initEventHandler(v1)
		}
	}

}

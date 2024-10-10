package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koma2211/ayan-capital_task/internal/entities"
)

func (h *Handler) initEventHandler(api *gin.RouterGroup) {
	events := api.Group("/events")
	{
		events.POST("/", h.addEvent())
	}
}

func (h *Handler) addEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req []entities.Event

		if err := c.BindJSON(&req); err != nil {
			response(c, http.StatusBadRequest, err.Error(), nil)
			return 
		}

		err := h.services.Eventer.AddEvents(c.Request.Context(), req)
		if err != nil {
			response(c, http.StatusInternalServerError, err.Error(), nil)
			return 
		}

		response(c, http.StatusOK, "success", nil)
	}
}
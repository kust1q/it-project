package handler

import "github.com/gin-gonic/gin"

func (h *Handler) info(c *gin.Context) {
	restaurant, err := h.showcase.Restaurant(c.Request.Context())
	h.respond(c, restaurant, err)
}

func (h *Handler) home(c *gin.Context) {
	home, err := h.showcase.Home(c.Request.Context())
	h.respond(c, home, err)
}

package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(menu MenuService, showcase ShowcaseService, logger *slog.Logger) *Handler {
	return &Handler{menu: menu, showcase: showcase, logger: logger}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })
	api := router.Group("/api/v1")
	api.GET("/categories", h.categories)
	api.GET("/categories/:id", h.category)
	api.GET("/dishes", h.dishes)
	api.GET("/dishes/:id", h.dish)
	api.GET("/info", h.info)
	api.GET("/home", h.home)
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "not_found", "message": "Endpoint not found"}})
	})
	return router
}

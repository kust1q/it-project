package handler

import (
	"errors"
	"net/http"
	"project/internal/domain/errs"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) respond(c *gin.Context, value any, err error) {
	if err == nil {
		c.JSON(http.StatusOK, value)
		return
	}
	if errors.Is(err, errs.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "not_found", "message": "Resource not found"}})
		return
	}
	h.logger.ErrorContext(c.Request.Context(), "request failed", "error", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "internal_error", "message": "Internal server error"}})
}

func invalid(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_request", "message": "Invalid request parameters"}})
}

func idParam(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		invalid(c)
		return 0, false
	}
	return id, true
}

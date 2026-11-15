package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"project/internal/domain/entity"
)

func (h *Handler) categories(c *gin.Context) {
	categories, err := h.menu.Categories(c.Request.Context())
	h.respond(c, categories, err)
}

func (h *Handler) category(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	category, err := h.menu.Category(c.Request.Context(), id)
	h.respond(c, category, err)
}

func (h *Handler) dish(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	dish, err := h.menu.Dish(c.Request.Context(), id)
	h.respond(c, dish, err)
}

func (h *Handler) dishes(c *gin.Context) {
	filter := entity.DishFilter{Page: 1, Limit: 20, Search: c.Query("search"), Sort: c.Query("sort")}
	for key, dst := range map[string]*int{"page": &filter.Page, "limit": &filter.Limit} {
		if value, ok := c.GetQuery(key); ok {
			n, err := strconv.Atoi(value)
			if err != nil || n < 1 || (key == "limit" && n > 100) || (key == "page" && n > 1000000) {
				invalid(c)
				return
			}
			*dst = n
		}
	}
	if value, ok := c.GetQuery("category_id"); ok {
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil || n < 1 {
			invalid(c)
			return
		}
		filter.CategoryID = n
	}
	switch filter.Sort {
	case "", "name", "-name", "price", "-price", "newest":
	default:
		invalid(c)
		return
	}
	dishes, err := h.menu.Dishes(c.Request.Context(), filter)
	h.respond(c, dishes, err)
}

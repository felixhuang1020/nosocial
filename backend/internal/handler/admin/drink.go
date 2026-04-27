package admin

import (
	"nosocial/internal/model"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DrinkAdminHandler struct {
	drinkService *service.DrinkService
}

func NewDrinkAdminHandler(drinkService *service.DrinkService) *DrinkAdminHandler {
	return &DrinkAdminHandler{drinkService: drinkService}
}

func (h *DrinkAdminHandler) List(c *gin.Context) {
	categoryID := uint32(0)
	status := int8(-1)
	page, size := getPageSize(c)
	if cid := c.Query("category_id"); cid != "" {
		if v, err := strconv.ParseUint(cid, 10, 32); err == nil {
			categoryID = uint32(v)
		}
	}
	if s := c.Query("status"); s != "" {
		status = int8(s[0] - '0')
	}

	drinks, total, err := h.drinkService.GetDrinkList(categoryID, status, page, size)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"list":  drinks,
		"total": total,
	})
}

func (h *DrinkAdminHandler) Create(c *gin.Context) {
	var drink model.Drink
	if err := c.ShouldBindJSON(&drink); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.drinkService.CreateDrink(&drink); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, drink)
}

func (h *DrinkAdminHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var drink model.Drink
	if err := c.ShouldBindJSON(&drink); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	drink.ID = id
	if err := h.drinkService.UpdateDrink(&drink); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, drink)
}

func (h *DrinkAdminHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err = h.drinkService.DeleteDrink(id); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *DrinkAdminHandler) Categories(c *gin.Context) {
	categories, err := h.drinkService.GetCategories()
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, categories)
}

func (h *DrinkAdminHandler) CreateCategory(c *gin.Context) {
	var category model.DrinkCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.drinkService.CreateCategory(&category); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, category)
}

func (h *DrinkAdminHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var category model.DrinkCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	category.ID = uint32(id)
	if err := h.drinkService.UpdateCategory(&category); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, category)
}

func (h *DrinkAdminHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.drinkService.DeleteCategory(uint32(id)); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}

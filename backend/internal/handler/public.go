package handler

import (
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PublicHandler struct {
	bannerService *service.BannerService
	drinkService  *service.DrinkService
	tarotService  *service.TarotService
	systemService *service.SystemService
}

func NewPublicHandler(bannerService *service.BannerService, drinkService *service.DrinkService, tarotService *service.TarotService, systemService *service.SystemService) *PublicHandler {
	return &PublicHandler{
		bannerService: bannerService,
		drinkService:  drinkService,
		tarotService:  tarotService,
		systemService: systemService,
	}
}

func (h *PublicHandler) ListBanners(c *gin.Context) {
	position := int8(1)
	if p := c.Query("position"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			position = int8(v)
		}
	}
	banners, err := h.bannerService.GetBannersByPosition(position)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, banners)
}

func (h *PublicHandler) ListDrinks(c *gin.Context) {
	categoryID := uint32(0)
	if cid := c.Query("category_id"); cid != "" {
		if v, err := strconv.ParseUint(cid, 10, 32); err == nil {
			categoryID = uint32(v)
		}
	}
	status := int8(1)
	page := 1
	size := 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			page = v
		}
	}
	if s := c.Query("size"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			size = v
		}
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

func (h *PublicHandler) GetDrinkDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	drink, err := h.drinkService.GetDrinkDetail(id)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, drink)
}

func (h *PublicHandler) ListTarotCards(c *gin.Context) {
	cards, err := h.tarotService.GetCards()
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, cards)
}

func (h *PublicHandler) GetConfig(c *gin.Context) {
	config := h.systemService.GetConfig()
	response.Success(c, config)
}

func (h *PublicHandler) GetShopInfo(c *gin.Context) {
	shop := h.systemService.GetShopInfo()
	response.Success(c, shop)
}

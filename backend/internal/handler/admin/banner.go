package admin

import (
	"nosocial/internal/model"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BannerAdminHandler struct {
	bannerService *service.BannerService
}

func NewBannerAdminHandler(bannerService *service.BannerService) *BannerAdminHandler {
	return &BannerAdminHandler{bannerService: bannerService}
}

func (h *BannerAdminHandler) List(c *gin.Context) {
	page, size := getPageSize(c)
	banners, total, err := h.bannerService.GetBannerList(page, size)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"list":  banners,
		"total": total,
	})
}

func (h *BannerAdminHandler) Create(c *gin.Context) {
	var banner model.Banner
	if err := c.ShouldBindJSON(&banner); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.bannerService.CreateBanner(&banner); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, banner)
}

func (h *BannerAdminHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var banner model.Banner
	if err := c.ShouldBindJSON(&banner); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	banner.ID = uint32(id)
	if err := h.bannerService.UpdateBanner(&banner); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, banner)
}

func (h *BannerAdminHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err = h.bannerService.DeleteBanner(uint32(id)); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}

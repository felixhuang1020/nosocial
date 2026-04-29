package admin

import (
	"errors"
	"fmt"
	"nosocial/config"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TarotAdminHandler struct {
	tarotService *service.TarotService
}

func NewTarotAdminHandler(tarotService *service.TarotService) *TarotAdminHandler {
	return &TarotAdminHandler{tarotService: tarotService}
}

func (h *TarotAdminHandler) Mappings(c *gin.Context) {
	mappings, err := h.tarotService.GetAllMappings()
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, mappings)
}

type updateMappingReq struct {
	Mappings []service.MappingItem `json:"mappings" binding:"required"`
}

func (h *TarotAdminHandler) UpdateMapping(c *gin.Context) {
	var req updateMappingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.tarotService.UpdateMappings(req.Mappings); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *TarotAdminHandler) ListCards(c *gin.Context) {
	cards, err := h.tarotService.GetCards()
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, cards)
}

type updateCardReq struct {
	ImageURL string `json:"image_url" binding:"required"`
}

// expectedOSSHostPrefix 返回本项目 OSS 直传域名前缀，用于校验上传后 URL 合法来源
func expectedOSSHostPrefix() string {
	cfg := config.C.OSS
	if cfg.Bucket == "" || cfg.Endpoint == "" {
		return ""
	}
	return fmt.Sprintf("https://%s.%s/", cfg.Bucket, cfg.Endpoint)
}

func (h *TarotAdminHandler) UpdateCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req updateCardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	url := strings.TrimSpace(req.ImageURL)
	if url == "" {
		response.BadRequest(c, "图片URL不能为空")
		return
	}
	if len(url) > 255 {
		response.BadRequest(c, "图片URL超长")
		return
	}
	// 仅允许本项目 OSS 域名，防注入外域链接
	if prefix := expectedOSSHostPrefix(); prefix != "" && !strings.HasPrefix(url, prefix) {
		response.BadRequest(c, "图片URL必须来自本项目 OSS 域名")
		return
	}
	if err := h.tarotService.UpdateCard(uint32(id), url); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.BadRequest(c, "卡牌不存在")
			return
		}
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}

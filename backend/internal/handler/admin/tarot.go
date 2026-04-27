package admin

import (
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"

	"github.com/gin-gonic/gin"
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

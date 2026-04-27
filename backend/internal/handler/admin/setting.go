package admin

import (
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"

	"github.com/gin-gonic/gin"
)

type SettingAdminHandler struct {
	systemService *service.SystemService
}

func NewSettingAdminHandler(systemService *service.SystemService) *SettingAdminHandler {
	return &SettingAdminHandler{systemService: systemService}
}

// GetSettings GET /admin/settings - 获取所有系统配置
func (h *SettingAdminHandler) GetSettings(c *gin.Context) {
	settings := h.systemService.GetAllSettings()
	response.Success(c, settings)
}

// UpdateSettings PUT /admin/settings - 更新系统配置
func (h *SettingAdminHandler) UpdateSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.systemService.UpdateSettings(req); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}

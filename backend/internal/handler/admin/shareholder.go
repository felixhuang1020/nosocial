package admin

import (
	"nosocial/internal/model"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShareholderAdminHandler struct {
	shareholderService *service.ShareholderService
	userService        *service.UserService
}

func NewShareholderAdminHandler(shareholderService *service.ShareholderService, userService *service.UserService) *ShareholderAdminHandler {
	return &ShareholderAdminHandler{
		shareholderService: shareholderService,
		userService:        userService,
	}
}

func (h *ShareholderAdminHandler) List(c *gin.Context) {
	page, size := getPageSize(c)
	users, total, err := h.shareholderService.GetShareholderList((page-1)*size, size)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}

	// Build response with team counts
	type shareholderItem struct {
		*model.User
		TeamCount int64 `json:"team_count"`
	}
	items := make([]shareholderItem, len(users))
	for i, u := range users {
		count, _ := h.shareholderService.GetTeamCount(u.ID)
		items[i] = shareholderItem{User: u, TeamCount: count}
	}

	response.Success(c, map[string]interface{}{
		"list":  items,
		"total": total,
	})
}

func (h *ShareholderAdminHandler) Earnings(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	page, size := getPageSize(c)
	records, total, err := h.shareholderService.GetEarnings(id, (page-1)*size, size)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"list":  records,
		"total": total,
	})
}

func (h *ShareholderAdminHandler) Team(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	team, err := h.shareholderService.GetTeam(id)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, team)
}

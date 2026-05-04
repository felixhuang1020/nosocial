package admin

import (
	"nosocial/internal/bootstrap"
	"nosocial/internal/middleware"
	"nosocial/internal/model"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	search := c.Query("search")
	users, total, err := h.shareholderService.GetShareholderList((page-1)*size, size, search)
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
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "无效的ID")
		return
	}
	page, size := getPageSize(c)
	// 操作审计：记录管理员访问的股东 ID
	adminID := middleware.GetAdminID(c)
	if bootstrap.Log != nil {
		bootstrap.Log.Info("admin access shareholder earnings",
			zap.Uint32("admin_id", adminID),
			zap.Uint64("target_user_id", id),
		)
	}
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
	if err != nil || id == 0 {
		response.BadRequest(c, "无效的ID")
		return
	}
	adminID := middleware.GetAdminID(c)
	if bootstrap.Log != nil {
		bootstrap.Log.Info("admin access shareholder team",
			zap.Uint32("admin_id", adminID),
			zap.Uint64("target_user_id", id),
		)
	}
	team, err := h.shareholderService.GetTeam(id)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, team)
}

// ListWithdrawals 管理端提现列表
// query: status (可选 0/1/2，缺省-1 全部), page, size
func (h *ShareholderAdminHandler) ListWithdrawals(c *gin.Context) {
	status := int8(-1)
	if s := c.Query("status"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v >= 0 && v <= 2 {
			status = int8(v)
		}
	}
	page, size := getPageSize(c)
	list, total, err := h.shareholderService.ListAllWithdrawals(status, (page-1)*size, size)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"list":  list,
		"total": total,
	})
}

// ApproveWithdrawal 审核通过
func (h *ShareholderAdminHandler) ApproveWithdrawal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "无效的ID")
		return
	}
	adminID := middleware.GetAdminID(c)
	if bootstrap.Log != nil {
		bootstrap.Log.Info("admin approve withdrawal",
			zap.Uint32("admin_id", adminID),
			zap.Uint64("withdrawal_id", id),
		)
	}
	if err := h.shareholderService.ApproveWithdrawal(id); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}

// RejectWithdrawal 审核拒绝（退还余额）
func (h *ShareholderAdminHandler) RejectWithdrawal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)
	adminID := middleware.GetAdminID(c)
	if bootstrap.Log != nil {
		bootstrap.Log.Info("admin reject withdrawal",
			zap.Uint32("admin_id", adminID),
			zap.Uint64("withdrawal_id", id),
			zap.String("reason", req.Reason),
		)
	}
	if err := h.shareholderService.RejectWithdrawal(id, req.Reason); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}

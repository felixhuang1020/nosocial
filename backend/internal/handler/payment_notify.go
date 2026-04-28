package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"nosocial/internal/dao"
	"nosocial/internal/model"
	"nosocial/internal/pkg/logx"
	"nosocial/internal/pkg/wxpay"
	"nosocial/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PaymentNotifyHandler 微信支付 V3 统一回调处理
// 路由：POST /api/v1/payment/wxpay/notify
type PaymentNotifyHandler struct {
	wx                 *wxpay.Client
	orderService       *service.OrderService
	shareholderService *service.ShareholderService
	notifyDAO          *dao.PaymentNotifyLogDAO
	logger             *zap.Logger
}

func NewPaymentNotifyHandler(
	wx *wxpay.Client,
	orderService *service.OrderService,
	shareholderService *service.ShareholderService,
	notifyDAO *dao.PaymentNotifyLogDAO,
	logger *zap.Logger,
) *PaymentNotifyHandler {
	return &PaymentNotifyHandler{
		wx:                 wx,
		orderService:       orderService,
		shareholderService: shareholderService,
		notifyDAO:          notifyDAO,
		logger:             logger,
	}
}

// wxpayFail 按微信支付 V3 规范返回非 2xx + {code,message}
func wxpayFail(c *gin.Context, status int, code, msg string) {
	c.JSON(status, gin.H{"code": code, "message": msg})
}

// wxpayOK 按微信支付 V3 规范返回 2xx（通常是 200/204 都可）
func wxpayOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": ""})
}

// Notify 统一回调入口
func (h *PaymentNotifyHandler) Notify(c *gin.Context) {
	if !h.wx.IsConfigured() {
		h.logger.Error("收到支付回调但 wxpay 未配置")
		wxpayFail(c, http.StatusServiceUnavailable, "FAIL", "wxpay not configured")
		return
	}

	content, err := h.wx.ParseNotify(c.Request.Context(), c.Request)
	if err != nil {
		h.logger.Warn("支付回调验签/解密失败", zap.Error(err))
		wxpayFail(c, http.StatusBadRequest, "FAIL", "invalid notify")
		return
	}
	if content == nil || content.OutTradeNo == "" || content.TransactionID == "" {
		wxpayFail(c, http.StatusBadRequest, "FAIL", "empty notify content")
		return
	}
	// 仅处理成交（其它状态忽略即可，微信会视为接收成功）
	if content.TradeState != "SUCCESS" {
		h.logger.Info("支付回调非 SUCCESS 状态，忽略",
			zap.String("out_trade_no", content.OutTradeNo),
			zap.String("trade_state", content.TradeState),
		)
		wxpayOK(c)
		return
	}

	raw, _ := json.Marshal(content)
	rawStr := string(raw)

	// 回调去重（transaction_id 全局唯一）+ 审计流水
	logRec := &model.PaymentNotifyLog{
		OrderNo:       content.OutTradeNo,
		TransactionID: content.TransactionID,
		OrderType:     content.Attach,
		RawBody:       rawStr,
	}
	if err := h.notifyDAO.Insert(logRec); err != nil {
		// 唯一键冲突 = 重复通知，视为幂等成功
		if isDupKey(err) {
			h.logger.Info("回调重复（transaction_id 已入库），忽略",
				zap.String("out_trade_no", content.OutTradeNo),
				zap.String("transaction_id", logx.MaskTx(content.TransactionID)),
			)
			wxpayOK(c)
			return
		}
		h.logger.Error("写入回调流水失败", zap.Error(err))
		wxpayFail(c, http.StatusInternalServerError, "FAIL", "server error")
		return
	}

	// 按 attach 分流业务处理
	switch content.Attach {
	case "drink":
		order, changed, err := h.orderService.PayCallback(
			content.OutTradeNo, content.TransactionID, rawStr, content.Amount.PayerTotal,
		)
		if err != nil {
			h.logger.Error("酒水订单回调处理失败",
				zap.String("out_trade_no", content.OutTradeNo),
				zap.Error(err),
			)
			wxpayFail(c, http.StatusInternalServerError, "FAIL", "process failed")
			return
		}
		if changed && order != nil {
			// 真正成交才计算佣金（避免幂等回调重复结算）
			if err := h.shareholderService.CalculateCommission(order); err != nil {
				h.logger.Error("佣金计算失败",
					zap.String("out_trade_no", content.OutTradeNo),
					zap.Error(err),
				)
				// 不阻塞回调确认，后续可靠对账补偿
			}
		}
	case "shareholder":
		if _, err := h.shareholderService.PayCallback(
			content.OutTradeNo, content.TransactionID, rawStr, content.Amount.PayerTotal,
		); err != nil {
			h.logger.Error("股东订单回调处理失败",
				zap.String("out_trade_no", content.OutTradeNo),
				zap.Error(err),
			)
			wxpayFail(c, http.StatusInternalServerError, "FAIL", "process failed")
			return
		}
	default:
		h.logger.Warn("未知 attach，忽略处理",
			zap.String("out_trade_no", content.OutTradeNo),
			zap.String("attach", content.Attach),
		)
	}

	_ = h.notifyDAO.MarkProcessed(content.TransactionID)
	wxpayOK(c)
}

// isDupKey 简单判断是否是唯一键冲突
func isDupKey(err error) bool {
	if err == nil {
		return false
	}
	// GORM 在 PostgreSQL 下重复键会返回包装后的 error，字符串中一般包含 "duplicate key"
	msg := err.Error()
	return errors.Is(err, gorm.ErrDuplicatedKey) ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "uk_pnl_tx_id") ||
		strings.Contains(msg, "UNIQUE constraint")
}

package handler

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strings"

	"nosocial/config"
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

	// 幂等检查：若流水中已存在该 transaction_id，表明之前已成功处理过，直接返回 OK
	if exists, err := h.notifyDAO.ExistsTxID(content.TransactionID); err != nil {
		h.logger.Error("查询回调流水失败", zap.Error(err))
		wxpayFail(c, http.StatusInternalServerError, "FAIL", "server error")
		return
	} else if exists {
		h.logger.Info("回调重复（transaction_id 已处理），忽略",
			zap.String("out_trade_no", content.OutTradeNo),
			zap.String("transaction_id", logx.MaskTx(content.TransactionID)),
		)
		wxpayOK(c)
		return
	}

	// 按 attach 分流业务处理 —— 业务失败不写流水，以便微信重推重试
	switch content.Attach {
	case "drink":
		order, changed, err := h.orderService.PayCallback(
			content.OutTradeNo, content.TransactionID, rawStr, content.Amount.PayerTotal,
		)
		if err != nil {
			// 金额不匹配是业务异常不是系统错误，返 OK 让微信停止重推，仅告警
			if strings.Contains(err.Error(), "amount mismatch") {
				h.logger.Warn("酒水订单金额不一致，忽略重推",
					zap.String("out_trade_no", content.OutTradeNo),
					zap.Error(err),
				)
				wxpayOK(c)
				return
			}
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
		// 前置金额校验：若金额与当前配置/订单不符，视为异常事件，
		// 返回 OK 让微信停止重推（重推永远不会成功），并记录告警供人工排查。
		expectedFen := int64(math.Round(config.C.Business.ShareholderFee * 100))
		if content.Amount.PayerTotal > 0 && expectedFen > 0 && content.Amount.PayerTotal != expectedFen {
			h.logger.Warn("股东回调金额与配置不符，忽略处理（人工告警）",
				zap.String("out_trade_no", content.OutTradeNo),
				zap.Int64("expect_fen", expectedFen),
				zap.Int64("got_fen", content.Amount.PayerTotal),
			)
			wxpayOK(c)
			return
		}
		if _, err := h.shareholderService.PayCallback(
			content.OutTradeNo, content.TransactionID, rawStr, content.Amount.PayerTotal,
		); err != nil {
			// 金额不匹配是业务异常不是系统错误，不应让微信重推
			if strings.Contains(err.Error(), "amount mismatch") {
				h.logger.Warn("股东订单金额不一致，忽略重推",
					zap.String("out_trade_no", content.OutTradeNo),
					zap.Error(err),
				)
				wxpayOK(c)
				return
			}
			h.logger.Error("股东订单回调处理失败",
				zap.String("out_trade_no", content.OutTradeNo),
				zap.Error(err),
			)
			wxpayFail(c, http.StatusInternalServerError, "FAIL", "process failed")
			return
		}
	default:
		// 未知 attach 视为异常事件：不写流水，记录 Error 供人工排查，
		// 同时返 OK 避免微信无限重推（重推仍会到达这里）。
		h.logger.Error("未知 attach，拒绝处理",
			zap.String("out_trade_no", content.OutTradeNo),
			zap.String("attach", content.Attach),
		)
		wxpayOK(c)
		return
	}

	// 业务成功后再写入流水（transaction_id 唯一索引保证后续重推幂等）
	logRec := &model.PaymentNotifyLog{
		OrderNo:       content.OutTradeNo,
		TransactionID: content.TransactionID,
		OrderType:     content.Attach,
		RawBody:       rawStr,
		Processed:     1,
	}
	if err := h.notifyDAO.Insert(logRec); err != nil && !isDupKey(err) {
		h.logger.Error("写入回调流水失败（业务已处理）", zap.Error(err))
		// 业务已处理成功，流水写入失败不再让微信重推，返回 OK
	}
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

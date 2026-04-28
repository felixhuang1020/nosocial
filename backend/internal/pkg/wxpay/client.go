// Package wxpay 封装微信支付 V3 客户端（JSAPI 下单 + 回调验签 + 查询）
package wxpay

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// ErrNotConfigured 支付未配置完成
var ErrNotConfigured = errors.New("wxpay: client not configured")

// Config 初始化所需参数
type Config struct {
	AppID             string
	MchID             string
	MchCertSerial     string
	MchPrivateKeyPath string
	APIV3Key          string
	NotifyURL         string
}

// Client 微信支付 V3 封装客户端
type Client struct {
	cfg        Config
	core       *core.Client
	jsapi      *jsapi.JsapiApiService
	handler    *notify.Handler
	privateKey *rsa.PrivateKey
}

// NewClient 构造微信支付客户端
// 当 MchID / MchCertSerial / MchPrivateKeyPath / APIV3Key 任一缺失时返回 nil, nil，
// 调用方应视为未配置，业务侧返回 503/配置错误，避免服务启动被阻塞。
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.MchID == "" || cfg.MchCertSerial == "" || cfg.MchPrivateKeyPath == "" || cfg.APIV3Key == "" || cfg.AppID == "" {
		return nil, nil
	}

	privateKey, err := utils.LoadPrivateKeyWithPath(cfg.MchPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load mch private key: %w", err)
	}

	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(cfg.MchID, cfg.MchCertSerial, privateKey, cfg.APIV3Key),
	}
	coreCli, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("new wxpay core client: %w", err)
	}

	// 注册证书下载器（notify 验签使用）
	mgr := downloader.MgrInstance()
	if err := mgr.RegisterDownloaderWithPrivateKey(ctx, privateKey, cfg.MchCertSerial, cfg.MchID, cfg.APIV3Key); err != nil {
		return nil, fmt.Errorf("register cert downloader: %w", err)
	}

	certVisitor := mgr.GetCertificateVisitor(cfg.MchID)
	handler, err := notify.NewRSANotifyHandler(cfg.APIV3Key, verifiers.NewSHA256WithRSAVerifier(certVisitor))
	if err != nil {
		return nil, fmt.Errorf("new notify handler: %w", err)
	}

	return &Client{
		cfg:        cfg,
		core:       coreCli,
		jsapi:      &jsapi.JsapiApiService{Client: coreCli},
		handler:    handler,
		privateKey: privateKey,
	}, nil
}

// IsConfigured 是否配置完成
func (c *Client) IsConfigured() bool {
	return c != nil && c.core != nil && c.jsapi != nil
}

// PrepayParams 小程序/JSAPI 端调起支付所需参数
type PrepayParams struct {
	PrepayID  string `json:"prepay_id"`
	AppID     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

// PrepayJSAPIReq 统一下单参数
type PrepayJSAPIReq struct {
	OutTradeNo  string // 商户订单号
	Description string // 商品描述
	AmountFen   int64  // 订单金额（分）
	OpenID      string // 支付用户 OpenID
	Attach      string // 附加数据，回调原样透传（建议存 order_type）
}

// PrepayJSAPI 预下单并返回小程序调起支付的参数
func (c *Client) PrepayJSAPI(ctx context.Context, req PrepayJSAPIReq) (*PrepayParams, error) {
	if !c.IsConfigured() {
		return nil, ErrNotConfigured
	}
	if req.OutTradeNo == "" || req.OpenID == "" || req.AmountFen <= 0 {
		return nil, fmt.Errorf("invalid prepay request")
	}

	resp, _, err := c.jsapi.PrepayWithRequestPayment(ctx, jsapi.PrepayRequest{
		Appid:       core.String(c.cfg.AppID),
		Mchid:       core.String(c.cfg.MchID),
		Description: core.String(req.Description),
		OutTradeNo:  core.String(req.OutTradeNo),
		NotifyUrl:   core.String(c.cfg.NotifyURL),
		Attach:      core.String(req.Attach),
		Amount: &jsapi.Amount{
			Currency: core.String("CNY"),
			Total:    core.Int64(req.AmountFen),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(req.OpenID),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("jsapi prepay: %w", err)
	}
	if resp == nil || resp.PrepayId == nil {
		return nil, fmt.Errorf("jsapi prepay: empty response")
	}

	return &PrepayParams{
		PrepayID:  derefString(resp.PrepayId),
		AppID:     derefString(resp.Appid),
		TimeStamp: derefString(resp.TimeStamp),
		NonceStr:  derefString(resp.NonceStr),
		Package:   derefString(resp.Package),
		SignType:  derefString(resp.SignType),
		PaySign:   derefString(resp.PaySign),
	}, nil
}

// BuildMiniPaySign 仅根据已有的 prepay_id 重新生成小程序 wx.requestPayment 参数，
// 用于订单已下单但前端丢失参数时的二次拉起（避免重复下单）。
func (c *Client) BuildMiniPaySign(ctx context.Context, prepayID string) (*PrepayParams, error) {
	if !c.IsConfigured() {
		return nil, ErrNotConfigured
	}
	if prepayID == "" {
		return nil, fmt.Errorf("empty prepay_id")
	}
	nonce, err := utils.GenerateNonce()
	if err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	pkg := "prepay_id=" + prepayID
	message := fmt.Sprintf("%s\n%s\n%s\n%s\n", c.cfg.AppID, ts, nonce, pkg)
	sig, err := c.core.Sign(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("sign: %w", err)
	}
	return &PrepayParams{
		PrepayID:  prepayID,
		AppID:     c.cfg.AppID,
		TimeStamp: ts,
		NonceStr:  nonce,
		Package:   pkg,
		SignType:  "RSA",
		PaySign:   sig.Signature,
	}, nil
}

// NotifyContent 微信支付 V3 回调解密后的结构（payments/transactions）
// 仅取常用字段，其它字段以透传原文方式落库
type NotifyContent struct {
	AppID          string `json:"appid"`
	MchID          string `json:"mchid"`
	OutTradeNo     string `json:"out_trade_no"`
	TransactionID  string `json:"transaction_id"`
	TradeType      string `json:"trade_type"`
	TradeState     string `json:"trade_state"`
	TradeStateDesc string `json:"trade_state_desc"`
	BankType       string `json:"bank_type"`
	Attach         string `json:"attach"`
	SuccessTime    string `json:"success_time"`
	Payer          struct {
		OpenID string `json:"openid"`
	} `json:"payer"`
	Amount struct {
		Total         int64  `json:"total"`
		PayerTotal    int64  `json:"payer_total"`
		Currency      string `json:"currency"`
		PayerCurrency string `json:"payer_currency"`
	} `json:"amount"`
}

// ParseNotify 验签并解密回调请求体，返回解密后结构体
func (c *Client) ParseNotify(ctx context.Context, r *http.Request) (*NotifyContent, error) {
	if !c.IsConfigured() {
		return nil, ErrNotConfigured
	}
	var content NotifyContent
	if _, err := c.handler.ParseNotifyRequest(ctx, r, &content); err != nil {
		return nil, err
	}
	return &content, nil
}

// QueryByOutTradeNo 以商户单号查询订单（对账 / 未收到回调时兜底查询）
func (c *Client) QueryByOutTradeNo(ctx context.Context, outTradeNo string) (*NotifyContent, error) {
	if !c.IsConfigured() {
		return nil, ErrNotConfigured
	}
	resp, _, err := c.jsapi.QueryOrderByOutTradeNo(ctx, jsapi.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(outTradeNo),
		Mchid:      core.String(c.cfg.MchID),
	})
	if err != nil {
		return nil, fmt.Errorf("query order: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("query order: empty response")
	}
	nc := &NotifyContent{
		AppID:          derefString(resp.Appid),
		MchID:          derefString(resp.Mchid),
		OutTradeNo:     derefString(resp.OutTradeNo),
		TransactionID:  derefString(resp.TransactionId),
		TradeType:      derefString(resp.TradeType),
		TradeState:     derefString(resp.TradeState),
		TradeStateDesc: derefString(resp.TradeStateDesc),
		BankType:       derefString(resp.BankType),
		Attach:         derefString(resp.Attach),
		SuccessTime:    derefString(resp.SuccessTime),
	}
	if resp.Payer != nil {
		nc.Payer.OpenID = derefString(resp.Payer.Openid)
	}
	if resp.Amount != nil {
		if resp.Amount.Total != nil {
			nc.Amount.Total = int64(*resp.Amount.Total)
		}
		if resp.Amount.PayerTotal != nil {
			nc.Amount.PayerTotal = int64(*resp.Amount.PayerTotal)
		}
		nc.Amount.Currency = derefString(resp.Amount.Currency)
		nc.Amount.PayerCurrency = derefString(resp.Amount.PayerCurrency)
	}
	return nc, nil
}

func derefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

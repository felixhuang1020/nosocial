package bootstrap

import (
	"context"

	"nosocial/config"
	"nosocial/internal/pkg/logx"
	"nosocial/internal/pkg/wxpay"

	"go.uber.org/zap"
)

// InitWXPay 初始化微信支付 V3 客户端
// 当配置缺失时返回 nil 客户端（允许未接入支付的环境启动），业务侧按未配置处理。
func InitWXPay(cfg *config.WXConfig, logger *zap.Logger) (*wxpay.Client, error) {
	c, err := wxpay.NewClient(context.Background(), wxpay.Config{
		AppID:             cfg.AppID,
		MchID:             cfg.MchID,
		MchCertSerial:     cfg.MchCertSerial,
		MchPrivateKeyPath: cfg.MchPrivateKeyPath,
		APIV3Key:          cfg.APIV3Key,
		NotifyURL:         cfg.NotifyURL,
	})
	if err != nil {
		return nil, err
	}
	if c == nil || !c.IsConfigured() {
		logger.Warn("wxpay 未配置完成，支付接口将返回未配置错误",
			zap.String("mch_id", logx.MaskMid(cfg.MchID)),
			zap.Bool("has_appid", cfg.AppID != ""),
			zap.Bool("has_cert_serial", cfg.MchCertSerial != ""),
			zap.Bool("has_key_path", cfg.MchPrivateKeyPath != ""),
			zap.Bool("has_apiv3_key", cfg.APIV3Key != ""),
			zap.Bool("has_notify_url", cfg.NotifyURL != ""),
		)
		return c, nil
	}
	logger.Info("wxpay 客户端初始化完成",
		zap.String("mch_id", logx.MaskMid(cfg.MchID)),
		zap.String("notify_url", cfg.NotifyURL),
	)
	return c, nil
}

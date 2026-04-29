package bootstrap

import (
	"fmt"
	"nosocial/config"
	"nosocial/internal/pkg/wxpay"

	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	DB     *gorm.DB
	Redis  *redis.Client
	Logger *zap.Logger
	Casbin *casbin.Enforcer
	WXPay  *wxpay.Client
}

func NewApp(configPath string) (*App, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	// 校验关键安全配置
	if cfg.JWT.AdminSecret == "" {
		return nil, fmt.Errorf("jwt.admin_secret 不能为空，请设置环境变量 NOSOCIAL_JWT_ADMIN_SECRET")
	}
	if cfg.JWT.WXSecret == "" {
		return nil, fmt.Errorf("jwt.wx_secret 不能为空，请设置环境变量 NOSOCIAL_JWT_WX_SECRET")
	}
	// 强度校验：release 模式下 secret 长度至少 32 字节（避免弱 JWT 签名密钥）
	if cfg.App.Mode == "release" {
		if len(cfg.JWT.AdminSecret) < 32 {
			return nil, fmt.Errorf("jwt.admin_secret 长度不足 32 字节，生产环境不允许弱密钥")
		}
		if len(cfg.JWT.WXSecret) < 32 {
			return nil, fmt.Errorf("jwt.wx_secret 长度不足 32 字节，生产环境不允许弱密钥")
		}
	}

	logger, err := InitLogger(&cfg.Log)
	if err != nil {
		return nil, err
	}

	db, err := InitDB(&cfg.PostgreSQL)
	if err != nil {
		logger.Error("init db failed", zap.Error(err))
		return nil, err
	}

	redisClient, err := InitRedis(&cfg.Redis)
	if err != nil {
		logger.Error("init redis failed", zap.Error(err))
		return nil, fmt.Errorf("init redis failed: %w", err)
	}

	enforcer, err := InitCasbin(db)
	if err != nil {
		logger.Error("init casbin failed", zap.Error(err))
		return nil, fmt.Errorf("init casbin failed: %w", err)
	}

	// 微信支付为可选组件：release 模式下未配置完整会 fail-fast；debug 模式仅警告
	wxpayCli, err := InitWXPay(&cfg.WX, logger)
	if err != nil {
		logger.Error("init wxpay failed", zap.Error(err))
		if cfg.App.Mode == "release" {
			return nil, fmt.Errorf("init wxpay failed: %w", err)
		}
	}

	return &App{
		Config: cfg,
		DB:     db,
		Redis:  redisClient,
		Logger: logger,
		Casbin: enforcer,
		WXPay:  wxpayCli,
	}, nil
}

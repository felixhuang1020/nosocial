package bootstrap

import (
	"fmt"
	"nosocial/config"

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
	}

	enforcer, err := InitCasbin(db)
	if err != nil {
		logger.Error("init casbin failed", zap.Error(err))
	}

	return &App{
		Config: cfg,
		DB:     db,
		Redis:  redisClient,
		Logger: logger,
		Casbin: enforcer,
	}, nil
}

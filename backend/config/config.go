package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var C *Config

type Config struct {
	App        AppConfig      `mapstructure:"app"`
	PostgreSQL PostgresConfig `mapstructure:"postgresql"`
	Redis      RedisConfig    `mapstructure:"redis"`
	JWT        JWTConfig      `mapstructure:"jwt"`
	OSS        OSSConfig      `mapstructure:"oss"`
	WX         WXConfig       `mapstructure:"wx"`
	Business   BusinessConfig `mapstructure:"business"`
	Log        LogConfig      `mapstructure:"log"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Mode    string `mapstructure:"mode"`
	Port    int    `mapstructure:"port"`
}

type PostgresConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	TimeZone        string `mapstructure:"timezone"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		p.Host, p.User, p.Password, p.DBName, p.Port, p.SSLMode, p.TimeZone)
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type JWTConfig struct {
	WXSecret    string `mapstructure:"wx_secret"`
	WXExpire    int    `mapstructure:"wx_expire"`
	AdminSecret string `mapstructure:"admin_secret"`
	AdminExpire int    `mapstructure:"admin_expire"`
}

type OSSConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	Bucket          string `mapstructure:"bucket"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	StsRoleArn      string `mapstructure:"sts_role_arn"`
	Region          string `mapstructure:"region"`
}

type WXConfig struct {
	AppID    string `mapstructure:"appid"`
	Secret   string `mapstructure:"secret"`
	MchID    string `mapstructure:"mch_id"`
	APIV3Key string `mapstructure:"api_v3_key"`
}

type BusinessConfig struct {
	ShareholderFee        float64 `mapstructure:"shareholder_fee"`
	CommissionRate        float64 `mapstructure:"commission_rate"`
	FreeDrinkID           int     `mapstructure:"free_drink_id"`
	ReviewCouponAmount    float64 `mapstructure:"review_coupon_amount"`
	ReviewCouponMinOrder  float64 `mapstructure:"review_coupon_min_order"`
	ReviewCouponValidDays int     `mapstructure:"review_coupon_valid_days"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// 启用环境变量覆盖，前缀为 NOSOCIAL_
	v.SetEnvPrefix("NOSOCIAL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// 环境变量兜底覆盖（NOSOCIAL_<SECTION>_<KEY>，应对 viper 嵌套结构 AutomaticEnv 不稳定的情况）
	// App
	if val := os.Getenv("NOSOCIAL_APP_MODE"); val != "" {
		cfg.App.Mode = val
	}
	if val := os.Getenv("NOSOCIAL_APP_PORT"); val != "" {
		if p, err := strconv.Atoi(val); err == nil {
			cfg.App.Port = p
		}
	}
	// PostgreSQL
	if val := os.Getenv("NOSOCIAL_POSTGRESQL_HOST"); val != "" {
		cfg.PostgreSQL.Host = val
	}
	if val := os.Getenv("NOSOCIAL_POSTGRESQL_PORT"); val != "" {
		if p, err := strconv.Atoi(val); err == nil {
			cfg.PostgreSQL.Port = p
		}
	}
	if val := os.Getenv("NOSOCIAL_POSTGRESQL_USER"); val != "" {
		cfg.PostgreSQL.User = val
	}
	if val := os.Getenv("NOSOCIAL_POSTGRESQL_PASSWORD"); val != "" {
		cfg.PostgreSQL.Password = val
	}
	if val := os.Getenv("NOSOCIAL_POSTGRESQL_DBNAME"); val != "" {
		cfg.PostgreSQL.DBName = val
	}
	// Redis
	if val := os.Getenv("NOSOCIAL_REDIS_HOST"); val != "" {
		cfg.Redis.Host = val
	}
	if val := os.Getenv("NOSOCIAL_REDIS_PORT"); val != "" {
		if p, err := strconv.Atoi(val); err == nil {
			cfg.Redis.Port = p
		}
	}
	if val := os.Getenv("NOSOCIAL_REDIS_PASSWORD"); val != "" {
		cfg.Redis.Password = val
	}
	if val := os.Getenv("NOSOCIAL_REDIS_DB"); val != "" {
		if p, err := strconv.Atoi(val); err == nil {
			cfg.Redis.DB = p
		}
	}
	// JWT
	if val := os.Getenv("NOSOCIAL_JWT_WX_SECRET"); val != "" {
		cfg.JWT.WXSecret = val
	}
	if val := os.Getenv("NOSOCIAL_JWT_ADMIN_SECRET"); val != "" {
		cfg.JWT.AdminSecret = val
	}
	// OSS
	if val := os.Getenv("NOSOCIAL_OSS_ACCESS_KEY_ID"); val != "" {
		cfg.OSS.AccessKeyID = val
	}
	if val := os.Getenv("NOSOCIAL_OSS_ACCESS_KEY_SECRET"); val != "" {
		cfg.OSS.AccessKeySecret = val
	}
	// WX
	if val := os.Getenv("NOSOCIAL_WX_APPID"); val != "" {
		cfg.WX.AppID = val
	}
	if val := os.Getenv("NOSOCIAL_WX_SECRET"); val != "" {
		cfg.WX.Secret = val
	}
	if val := os.Getenv("NOSOCIAL_WX_API_V3_KEY"); val != "" {
		cfg.WX.APIV3Key = val
	}

	C = &cfg
	return &cfg, nil
}

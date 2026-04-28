package bootstrap

import (
	"fmt"
	"nosocial/config"
	"nosocial/internal/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.PostgresConfig) (*gorm.DB, error) {
	logLevel := logger.Info
	if config.C.App.Mode == "release" {
		logLevel = logger.Silent
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	DB = db
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.ShareholderOrder{},
		&model.CommissionRecord{},
		&model.TarotCard{},
		&model.TarotReading{},
		&model.TarotDrinkMapping{},
		&model.DrinkCategory{},
		&model.Drink{},
		&model.DrinkOrder{},
		&model.BirthdayGift{},
		&model.Review{},
		&model.Coupon{},
		&model.Banner{},
		&model.AdminUser{},
		&model.Setting{},
		&model.Withdrawal{},
		&model.PaymentNotifyLog{},
	)
}

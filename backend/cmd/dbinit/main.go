package main

// 远程 PostgreSQL 初始化工具
// 用法: go run ./cmd/dbinit
// 功能: 1) 创建 nosocial 数据库 2) AutoMigrate 所有表 3) 植入种子数据

import (
	"fmt"
	"log"
	"nosocial/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	pgHost = "123.57.241.163"
	pgPort = 5432
	pgUser = "postgres"
	pgPwd  = "Nzi2001oak*."
	pgDB   = "nosocial"
)

func main() {
	// 1) 连接 postgres 系统库，创建业务库
	sysDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		pgHost, pgPort, pgUser, pgPwd)

	sysDB, err := gorm.Open(postgres.Open(sysDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("连接 postgres 系统库失败: %v", err)
	}
	log.Println("[1/4] 已连接到远程 PG postgres 库")

	var cnt int64
	sysDB.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", pgDB).Scan(&cnt)
	if cnt == 0 {
		if err := sysDB.Exec(fmt.Sprintf("CREATE DATABASE %s", pgDB)).Error; err != nil {
			log.Fatalf("创建数据库 %s 失败: %v", pgDB, err)
		}
		log.Printf("[2/4] 已创建数据库: %s", pgDB)
	} else {
		log.Printf("[2/4] 数据库 %s 已存在，跳过创建", pgDB)
	}

	// 2) 连接业务库，执行 AutoMigrate
	bizDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		pgHost, pgPort, pgUser, pgPwd, pgDB)
	db, err := gorm.Open(postgres.Open(bizDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
	if err != nil {
		log.Fatalf("连接业务库失败: %v", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.ShareholderOrder{},
		&model.CommissionRecord{},
		&model.Withdrawal{},
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
	); err != nil {
		log.Fatalf("AutoMigrate 失败: %v", err)
	}
	log.Println("[3/4] 所有表结构已同步")

	// 3) 种子数据
	seedAdmin(db)
	seedCategories(db)
	seedSettings(db)

	log.Println("[4/4] 种子数据植入完成")
	log.Println("========================================")
	log.Println(" 远程 PG 初始化完成")
	log.Println("   库名: nosocial")
	log.Println("   admin / admin123")
	log.Println("========================================")
}

func seedAdmin(db *gorm.DB) {
	var cnt int64
	db.Model(&model.AdminUser{}).Where("username = ?", "admin").Count(&cnt)
	if cnt > 0 {
		log.Println("  - admin 账号已存在")
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	nickname := "超级管理员"
	admin := &model.AdminUser{
		Username: "admin",
		Password: string(hash),
		Nickname: &nickname,
		Role:     3,
		Status:   1,
	}
	if err := db.Create(admin).Error; err != nil {
		log.Printf("  - 创建 admin 失败: %v", err)
		return
	}
	log.Println("  - 已创建 admin / admin123")
}

func seedCategories(db *gorm.DB) {
	names := []struct {
		Name string
		Sort int
	}{
		{"鸡尾酒", 1},
		{"啤酒", 2},
		{"烈酒", 3},
		{"无酒精", 4},
	}
	for _, c := range names {
		var cnt int64
		db.Model(&model.DrinkCategory{}).Where("name = ?", c.Name).Count(&cnt)
		if cnt > 0 {
			continue
		}
		db.Create(&model.DrinkCategory{Name: c.Name, SortOrder: c.Sort, Status: 1})
	}
	log.Println("  - 酒水分类: 4 个已就绪")
}

func seedSettings(db *gorm.DB) {
	kvs := []struct {
		Key, Val, Desc string
	}{
		{"shop_name", "NoSocial 酒吧", "店铺名称"},
		{"shop_address", "北京市朝阳区三里屯太古里北区 N8-20", "店铺地址"},
		{"business_hours", "周一至周日 18:00 - 04:00", "营业时间"},
		{"phone", "010-8888-6666", "联系电话"},
		{"latitude", "39.934", "纬度"},
		{"longitude", "116.455", "经度"},
		{"wifi_name", "NoSocial_Free", "WiFi名称"},
		{"wifi_password", "nosocial888", "WiFi密码"},
	}
	for _, kv := range kvs {
		var cnt int64
		db.Model(&model.Setting{}).Where("key = ?", kv.Key).Count(&cnt)
		if cnt > 0 {
			continue
		}
		desc := kv.Desc
		db.Create(&model.Setting{Key: kv.Key, Value: kv.Val, Desc: &desc})
	}
	log.Println("  - 系统配置: 8 条已就绪")
}

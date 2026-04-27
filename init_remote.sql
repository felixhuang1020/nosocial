-- ============================================================
-- NoSocial 数据库初始化脚本
-- 适用于 PostgreSQL 17+
-- 运行方式: psql -h <host> -U postgres -d postgres -f init.sql
-- ============================================================

-- 创建数据库（如果不存在）
SELECT 'CREATE DATABASE nosocial' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'nosocial')\gexec

\c nosocial;

-- ============================================================
-- 用户表
-- ============================================================
CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "openid" VARCHAR(64) UNIQUE NOT NULL,
    "unionid" VARCHAR(64),
    "nickname" VARCHAR(64),
    "avatar_url" VARCHAR(255),
    "gender" SMALLINT DEFAULT 0,
    "phone" VARCHAR(20),
    "status" SMALLINT NOT NULL DEFAULT 1,
    "shareholder_status" SMALLINT NOT NULL DEFAULT 0,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "admin_users" (
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR(32) UNIQUE NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "nickname" VARCHAR(32),
    "role" SMALLINT NOT NULL DEFAULT 1,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "last_login_at" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 初始化超级管理员账号（密码: admin123）
INSERT INTO admin_users (username, password, nickname, role, status, created_at)
VALUES ('admin', '$2a$10$EJTAyZBUmUU.iOHXYXi7cOuw8jBT/T98xa0KXDSB/WWSrSY2d5bpy', '超级管理员', 3, 1, NOW())
ON CONFLICT (username) DO NOTHING;

-- ============================================================
-- 酒水分类表
-- ============================================================
CREATE TABLE IF NOT EXISTS "drink_categories" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(32) NOT NULL,
    "sort_order" SMALLINT NOT NULL DEFAULT 0,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO drink_categories (name, sort_order, status) VALUES
    ('鸡尾酒', 1, 1),
    ('啤酒', 2, 1),
    ('烈酒', 3, 1),
    ('无酒精', 4, 1)
ON CONFLICT DO NOTHING;

-- ============================================================
-- 酒水表
-- ============================================================
CREATE TABLE IF NOT EXISTS "drinks" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(64) NOT NULL,
    "english_name" VARCHAR(64),
    "category_id" SMALLINT NOT NULL,
    "price" DECIMAL(10,2) NOT NULL,
    "cost_price" DECIMAL(10,2) DEFAULT 0,
    "description" TEXT,
    "image_url" VARCHAR(255),
    "alcohol_degree" VARCHAR(16),
    "is_recommended" SMALLINT NOT NULL DEFAULT 0,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- 订单表
-- ============================================================
CREATE TABLE IF NOT EXISTS "drink_orders" (
    "id" SERIAL PRIMARY KEY,
    "order_no" VARCHAR(32) UNIQUE NOT NULL,
    "user_id" INTEGER NOT NULL,
    "total_amount" DECIMAL(10,2) NOT NULL,
    "discount_amount" DECIMAL(10,2) NOT NULL DEFAULT 0,
    "pay_amount" DECIMAL(10,2) NOT NULL,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "pay_time" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "drink_order_items" (
    "id" SERIAL PRIMARY KEY,
    "order_id" INTEGER NOT NULL,
    "drink_id" INTEGER NOT NULL,
    "drink_name" VARCHAR(64) NOT NULL,
    "price" DECIMAL(10,2) NOT NULL,
    "quantity" INTEGER NOT NULL DEFAULT 1,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- 股东表
-- ============================================================
CREATE TABLE IF NOT EXISTS "shareholder_orders" (
    "id" SERIAL PRIMARY KEY,
    "order_no" VARCHAR(32) UNIQUE NOT NULL,
    "user_id" INTEGER NOT NULL,
    "level" SMALLINT NOT NULL DEFAULT 1,
    "amount" DECIMAL(10,2) NOT NULL,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "pay_time" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "commission_records" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "order_id" INTEGER,
    "amount" DECIMAL(10,2) NOT NULL,
    "type" SMALLINT NOT NULL DEFAULT 1,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- Tarot 塔罗牌表
-- ============================================================
CREATE TABLE IF NOT EXISTS "tarot_cards" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(64) NOT NULL,
    "english_name" VARCHAR(64),
    "description" TEXT,
    "image_url" VARCHAR(255),
    "interpretation" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "tarot_readings" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "cards_json" TEXT NOT NULL,
    "interpretation" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "tarot_drink_mappings" (
    "id" SERIAL PRIMARY KEY,
    "card_id" INTEGER NOT NULL,
    "drink_id" INTEGER NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- 营销表
-- ============================================================
CREATE TABLE IF NOT EXISTS "reviews" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "order_id" INTEGER,
    "rating" SMALLINT NOT NULL,
    "content" TEXT,
    "images" TEXT,
    "reply" TEXT,
    "status" SMALLINT NOT NULL DEFAULT 0,
    "audit_admin_id" INTEGER,
    "reject_reason" VARCHAR(255),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "coupons" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "code" VARCHAR(32) UNIQUE NOT NULL,
    "type" SMALLINT NOT NULL DEFAULT 1,
    "amount" DECIMAL(10,2) NOT NULL,
    "min_order_amount" DECIMAL(10,2) NOT NULL DEFAULT 0,
    "valid_from" TIMESTAMPTZ NOT NULL,
    "valid_until" TIMESTAMPTZ NOT NULL,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "used_at" TIMESTAMPTZ,
    "used_order_id" INTEGER,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "birthday_gifts" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "gift_type" SMALLINT NOT NULL DEFAULT 1,
    "status" SMALLINT NOT NULL DEFAULT 0,
    "claimed_at" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- Banner 表
-- ============================================================
CREATE TABLE IF NOT EXISTS "banners" (
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR(64) NOT NULL,
    "image_url" VARCHAR(255) NOT NULL,
    "link_type" SMALLINT NOT NULL DEFAULT 1,
    "link_value" VARCHAR(255),
    "sort_order" SMALLINT NOT NULL DEFAULT 0,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "start_time" TIMESTAMPTZ,
    "end_time" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- 系统配置表
-- ============================================================
CREATE TABLE IF NOT EXISTS "settings" (
    "id" SERIAL PRIMARY KEY,
    "key" VARCHAR(64) UNIQUE NOT NULL,
    "value" TEXT NOT NULL,
    "desc" VARCHAR(255),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO settings (key, value, desc) VALUES
    ('shop_name', 'NoSocial 酒吧', '店铺名称'),
    ('shop_address', '北京市朝阳区三里屯太古里北区 N8-20', '店铺地址'),
    ('business_hours', '周一至周日 18:00 - 04:00', '营业时间'),
    ('phone', '010-8888-6666', '联系电话'),
    ('latitude', '39.934', '纬度'),
    ('longitude', '116.455', '经度'),
    ('wifi_name', 'NoSocial_Free', 'WiFi名称'),
    ('wifi_password', 'nosocial888', 'WiFi密码')
ON CONFLICT (key) DO NOTHING;

-- ============================================================
-- 提现表
-- ============================================================
CREATE TABLE IF NOT EXISTS "withdrawals" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "amount" DECIMAL(10,2) NOT NULL,
    "status" SMALLINT NOT NULL DEFAULT 0,
    "audit_admin_id" INTEGER,
    "reject_reason" VARCHAR(255),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- 索引
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_users_openid ON users(openid);
CREATE INDEX IF NOT EXISTS idx_drinks_category ON drinks(category_id);
CREATE INDEX IF NOT EXISTS idx_drinks_status ON drinks(status);
CREATE INDEX IF NOT EXISTS idx_orders_user ON drink_orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON drink_orders(status);
CREATE INDEX IF NOT EXISTS idx_shareholder_user ON shareholder_orders(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_coupons_user ON coupons(user_id);
CREATE INDEX IF NOT EXISTS idx_coupons_status ON coupons(status);

-- ============================================================
-- 完成
-- ============================================================
SELECT 'NoSocial 数据库初始化完成!' AS result;

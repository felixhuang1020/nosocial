-- ============================================================
-- 初始化超级管理员账号
-- 安全策略：不再写入默认弱口令（旧 admin/admin123 已移除）
-- 部署时请在下方 INSERT 前自行生成 bcrypt 密文并替换占位符：
--   htpasswd -bnBC 12 '' 'YourStrongPasswordHere' | tr -d ':\n'
-- 或 Go：golang.org/x/crypto/bcrypt.GenerateFromPassword([]byte("..."), 12)
-- 然后取消注释执行。切勿把明文密码写入任何脚本或仓库。
-- ============================================================
-- INSERT INTO admin_users (username, password, nickname, role, status, created_at)
-- VALUES ('admin', '__BCRYPT_HASH_REPLACE_ME__', '超级管理员', 3, 1, NOW())
-- ON CONFLICT (username) DO NOTHING;

-- 初始化默认酒水分类
INSERT INTO drink_categories (name, sort_order, status) VALUES
('鸡尾酒', 1, 1),
('威士忌', 2, 1),
('精酿啤酒', 3, 1),
('葡萄酒', 4, 1),
('无酒精', 5, 1)
ON CONFLICT DO NOTHING;

-- 初始化默认酒水
INSERT INTO drinks (name, english_name, category_id, price, description, image_url, is_recommended, is_free_drink, status, sort_order) VALUES
('莫吉托', 'Mojito', 1, 68.00, '清爽的薄荷与青柠，朗姆酒的经典之作', '/images/mojito.jpg', 1, 0, 1, 1),
('长岛冰茶', 'Long Island Iced Tea', 1, 88.00, '五款基酒的经典调配，浓烈却不失风味', '/images/long-island.jpg', 1, 0, 1, 2),
('生日特调', 'Birthday Special', 1, 88.00, '专为寿星调制的特调鸡尾酒', '/images/birthday.jpg', 0, 1, 1, 3)
ON CONFLICT DO NOTHING;

-- ==================== 支付改造 DDL（幂等字段 + 回调日志）====================
ALTER TABLE IF EXISTS drink_orders ADD COLUMN IF NOT EXISTS prepay_id VARCHAR(64);
ALTER TABLE IF EXISTS drink_orders ADD COLUMN IF NOT EXISTS transaction_id VARCHAR(64);
ALTER TABLE IF EXISTS drink_orders ADD COLUMN IF NOT EXISTS items JSONB;
ALTER TABLE IF EXISTS drink_orders ADD COLUMN IF NOT EXISTS notify_raw JSONB;
CREATE UNIQUE INDEX IF NOT EXISTS uk_drink_tx_id ON drink_orders(transaction_id) WHERE transaction_id IS NOT NULL;

ALTER TABLE IF EXISTS shareholder_orders ADD COLUMN IF NOT EXISTS prepay_id VARCHAR(64);
ALTER TABLE IF EXISTS shareholder_orders ADD COLUMN IF NOT EXISTS notify_raw JSONB;

CREATE TABLE IF NOT EXISTS payment_notify_logs (
    id             BIGSERIAL PRIMARY KEY,
    order_no       VARCHAR(32)  NOT NULL,
    transaction_id VARCHAR(64)  NOT NULL,
    order_type     VARCHAR(8)   NOT NULL DEFAULT '',
    raw_body       JSONB,
    processed      SMALLINT     NOT NULL DEFAULT 0,
    received_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_pnl_tx_id ON payment_notify_logs(transaction_id);
CREATE INDEX IF NOT EXISTS idx_pnl_order_no ON payment_notify_logs(order_no);

-- ==================== 佣金幂等索引 ====================
-- 同一 (order_id, order_type) 组合只能存一条佣金，补偿重放不会翻倍发放
CREATE UNIQUE INDEX IF NOT EXISTS uk_commission_order ON commission_records(order_id, order_type);

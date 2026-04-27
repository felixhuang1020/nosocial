-- 初始化超级管理员账号（密码: admin123）
INSERT INTO admin_users (username, password, nickname, role, status, created_at)
VALUES ('admin', '$2a$10$EJTAyZBUmUU.iOHXYXi7cOuw8jBT/T98xa0KXDSB/WWSrSY2d5bpy', '超级管理员', 3, 1, NOW())
ON CONFLICT (username) DO NOTHING;

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

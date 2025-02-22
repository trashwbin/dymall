-- 创建管理员用户
INSERT INTO users (id, username, email, password, created_at, updated_at)
VALUES (1, 'admin', 'admin@example.com', '$2a$10$your_hashed_password', NOW(), NOW());

-- 创建测试用户
INSERT INTO users (id, username, email, password, created_at, updated_at)
VALUES (2, 'test_user', 'test@example.com', '$2a$10$your_hashed_password', NOW(), NOW());

-- 创建示例商品
INSERT INTO products (id, name, description, price, stock, created_at, updated_at)
VALUES 
(1, 'iPhone 14', '苹果最新旗舰手机', 5999.00, 100, NOW(), NOW()),
(2, 'MacBook Pro', '专业级笔记本电脑', 12999.00, 50, NOW(), NOW()),
(3, 'AirPods Pro', '主动降噪耳机', 1999.00, 200, NOW(), NOW());

-- 创建商品分类
INSERT INTO categories (id, name, created_at, updated_at)
VALUES 
(1, '手机', NOW(), NOW()),
(2, '电脑', NOW(), NOW()),
(3, '配件', NOW(), NOW());

-- 商品分类关联
INSERT INTO product_categories (product_id, category_id)
VALUES 
(1, 1),  -- iPhone 14 属于手机类
(2, 2),  -- MacBook Pro 属于电脑类
(3, 3);  -- AirPods Pro 属于配件类

-- 为管理员分配角色
INSERT INTO casbin_rule (ptype, v0, v1, v2)
VALUES ('g', 'user:1', 'admin', '-');

-- 为测试用户分配角色
INSERT INTO casbin_rule (ptype, v0, v1, v2)
VALUES ('g', 'user:2', 'user', '-'); 
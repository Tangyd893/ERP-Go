-- DEPRECATED: 归档快照，不参与数据库迁移。业务种子数据请使用 iam-service/migrations/002_seed_dev_admin.sql
-- 初始化默认超级管理员和基础数据
-- 密码为 "admin123"，bcrypt 哈希

INSERT INTO users (id, tenant_id, username, password_hash, nickname, status)
VALUES ('u-root-admin', 't-default', 'admin',
        '$2a$10$0bq47I5qJ5GHKlJewuj6jetY2hor.uYTseL4WLDITwIYhPo2pQole',
        '超级管理员', 'active')
ON CONFLICT (id) DO NOTHING;

INSERT INTO roles (id, tenant_id, name, code, description)
VALUES ('r-super-admin', 't-default', '超级管理员', 'super_admin', '系统默认超级管理员角色')
ON CONFLICT (id) DO NOTHING;

INSERT INTO roles (id, tenant_id, name, code, description)
VALUES ('r-warehouse-op', 't-default', '仓库操作员', 'warehouse_operator', '仓库作业人员角色')
ON CONFLICT (id) DO NOTHING;

INSERT INTO roles (id, tenant_id, name, code, description)
VALUES ('r-finance', 't-default', '财务人员', 'finance', '财务结算人员角色')
ON CONFLICT (id) DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
VALUES ('u-root-admin', 'r-super-admin')
ON CONFLICT (user_id, role_id) DO NOTHING;

-- 基础菜单和按钮权限
INSERT INTO permissions (id, code, name, resource_type, action, sort_order)
VALUES
  ('p-dashboard', 'dashboard', '首页看板', 'menu', 'read', 1),
  ('p-product-mgmt', 'product:menu', '商品管理', 'menu', 'read', 2),
  ('p-product-list', 'product:list', '商品列表', 'button', 'list', 3),
  ('p-channel-mgmt', 'channel:menu', '渠道管理', 'menu', 'read', 4),
  ('p-channel-store', 'channel:store_list', '店铺管理', 'button', 'list', 5),
  ('p-order-mgmt', 'order:menu', '订单管理', 'menu', 'read', 6),
  ('p-order-list', 'order:list', '订单列表', 'button', 'list', 7),
  ('p-inventory-mgmt', 'inventory:menu', '库存管理', 'menu', 'read', 8),
  ('p-inventory-list', 'inventory:list', '库存查询', 'button', 'list', 9),
  ('p-warehouse-mgmt', 'warehouse:menu', '仓储管理', 'menu', 'read', 10),
  ('p-warehouse-out', 'warehouse:outbound_list', '出库管理', 'button', 'list', 11),
  ('p-transport-mgmt', 'transport:menu', '物流管理', 'menu', 'read', 12),
  ('p-system-mgmt', 'system:menu', '系统管理', 'menu', 'read', 20),
  ('p-user-mgmt', 'user:menu', '用户管理', 'menu', 'read', 21),
  ('p-user-list', 'user:list', '用户列表', 'button', 'list', 22),
  ('p-user-create', 'user:create', '创建用户', 'button', 'create', 23),
  ('p-user-update', 'user:update', '编辑用户', 'button', 'update', 24),
  ('p-user-disable', 'user:disable', '禁用用户', 'button', 'delete', 25),
  ('p-user-enable', 'user:enable', '启用用户', 'button', 'update', 26),
  ('p-user-role', 'user:assign_role', '分配角色', 'button', 'update', 27),
  ('p-role-mgmt', 'role:menu', '角色管理', 'menu', 'read', 30),
  ('p-role-list', 'role:list', '角色列表', 'button', 'list', 31),
  ('p-role-create', 'role:create', '创建角色', 'button', 'create', 32),
  ('p-role-update', 'role:update', '编辑角色', 'button', 'update', 33),
  ('p-role-delete', 'role:delete', '删除角色', 'button', 'delete', 34),
  ('p-role-perm', 'role:assign_perm', '分配权限', 'button', 'update', 35),
  ('p-audit-mgmt', 'audit:menu', '操作审计', 'menu', 'read', 40),
  ('p-audit-list', 'audit:list', '审计日志', 'button', 'list', 41)
ON CONFLICT (code) DO NOTHING;

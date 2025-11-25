-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(128) UNIQUE,
    real_name VARCHAR(64),
    phone VARCHAR(32),
    avatar VARCHAR(255),
    status VARCHAR(32) DEFAULT 'active',
    is_admin BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_username ON users (username);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_status ON users (status);

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, role_id)
);

CREATE INDEX idx_user_roles_user_id ON user_roles (user_id);

CREATE INDEX idx_user_roles_role_id ON user_roles (role_id);

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE,
    resource VARCHAR(64) NOT NULL,
    action VARCHAR(32) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    id SERIAL PRIMARY KEY,
    role_id INTEGER NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role_id ON role_permissions (role_id);

CREATE INDEX idx_role_permissions_permission_id ON role_permissions (permission_id);

-- 插入默认角色
INSERT INTO
    roles (name, description)
VALUES ('admin', '系统管理员'),
    ('user', '普通用户')
ON CONFLICT (name) DO NOTHING;

-- 插入默认权限
INSERT INTO
    permissions (
        name,
        resource,
        action,
        description
    )
VALUES (
        'kb_create',
        'knowledge_base',
        'create',
        '创建知识库'
    ),
    (
        'kb_read',
        'knowledge_base',
        'read',
        '查看知识库'
    ),
    (
        'kb_update',
        'knowledge_base',
        'update',
        '更新知识库'
    ),
    (
        'kb_delete',
        'knowledge_base',
        'delete',
        '删除知识库'
    ),
    (
        'doc_create',
        'document',
        'create',
        '上传文档'
    ),
    (
        'doc_read',
        'document',
        'read',
        '查看文档'
    ),
    (
        'doc_delete',
        'document',
        'delete',
        '删除文档'
    ),
    (
        'chat_create',
        'chat',
        'create',
        '创建对话'
    ),
    (
        'chat_read',
        'chat',
        'read',
        '查看对话'
    ),
    (
        'user_manage',
        'user',
        'manage',
        '管理用户'
    )
ON CONFLICT (name) DO NOTHING;

-- 为管理员角色分配所有权限
INSERT INTO
    role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE
    r.name = 'admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 为普通用户角色分配基础权限
INSERT INTO
    role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE
    r.name = 'user'
    AND p.name IN (
        'kb_create',
        'kb_read',
        'kb_update',
        'doc_create',
        'doc_read',
        'chat_create',
        'chat_read'
    )
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 创建默认管理员用户（密码：123456，bcrypt 加密）
INSERT INTO
    users (
        username,
        password,
        email,
        real_name,
        is_admin,
        status
    )
VALUES (
        'admin',
        '$2a$10$PW0oux5IncrkVyyTRhnig.FkUleL5MpJo4BwpFU53FRiRxkf4t1fC',
        'admin@ragljx.com',
        '系统管理员',
        true,
        'active'
    )
ON CONFLICT (username) DO NOTHING;

-- 为默认管理员分配管理员角色
INSERT INTO
    user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE
    u.username = 'admin'
    AND r.name = 'admin'
ON CONFLICT (user_id, role_id) DO NOTHING;
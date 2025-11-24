-- 006_init_default_admin.sql
-- 创建默认管理员用户和权限

-- 插入默认权限
INSERT INTO
    permissions (
        id,
        name,
        code,
        resource,
        action,
        description,
        created_at,
        updated_at
    )
VALUES (
        'perm-001',
        '用户管理',
        'user:manage',
        'user',
        '*',
        '用户的增删改查',
        NOW(),
        NOW()
    ),
    (
        'perm-002',
        '角色管理',
        'role:manage',
        'role',
        '*',
        '角色的增删改查',
        NOW(),
        NOW()
    ),
    (
        'perm-003',
        '权限管理',
        'permission:manage',
        'permission',
        '*',
        '权限的增删改查',
        NOW(),
        NOW()
    ),
    (
        'perm-004',
        '知识库管理',
        'kb:manage',
        'knowledge_base',
        '*',
        '知识库的增删改查',
        NOW(),
        NOW()
    ),
    (
        'perm-005',
        '文档管理',
        'doc:manage',
        'document',
        '*',
        '文档的增删改查',
        NOW(),
        NOW()
    ),
    (
        'perm-006',
        '对话管理',
        'chat:manage',
        'chat',
        '*',
        '对话的增删改查',
        NOW(),
        NOW()
    ),
    (
        'perm-007',
        '系统配置',
        'system:config',
        'system',
        '*',
        '系统配置管理',
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- 插入默认角色
INSERT INTO
    roles (
        id,
        name,
        code,
        description,
        created_at,
        updated_at
    )
VALUES (
        'role-admin',
        '超级管理员',
        'admin',
        '拥有所有权限的超级管理员',
        NOW(),
        NOW()
    ),
    (
        'role-user',
        '普通用户',
        'user',
        '普通用户，可以使用知识库和对话功能',
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- 为超级管理员角色分配所有权限
INSERT INTO
    role_permissions (
        role_id,
        permission_id,
        created_at
    )
VALUES (
        'role-admin',
        'perm-001',
        NOW()
    ),
    (
        'role-admin',
        'perm-002',
        NOW()
    ),
    (
        'role-admin',
        'perm-003',
        NOW()
    ),
    (
        'role-admin',
        'perm-004',
        NOW()
    ),
    (
        'role-admin',
        'perm-005',
        NOW()
    ),
    (
        'role-admin',
        'perm-006',
        NOW()
    ),
    (
        'role-admin',
        'perm-007',
        NOW()
    )
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 为普通用户角色分配基本权限
INSERT INTO
    role_permissions (
        role_id,
        permission_id,
        created_at
    )
VALUES (
        'role-user',
        'perm-004',
        NOW()
    ),
    (
        'role-user',
        'perm-005',
        NOW()
    ),
    (
        'role-user',
        'perm-006',
        NOW()
    )
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 插入默认管理员用户
-- 密码: 123456 (bcrypt hash)
-- 使用 bcrypt cost=10 生成的 hash
INSERT INTO
    users (
        id,
        username,
        email,
        password,
        nickname,
        avatar,
        status,
        created_at,
        updated_at
    )
VALUES (
        'user-admin',
        'admin',
        'admin@ragljx.com',
        '$2a$10$PW0oux5IncrkVyyTRhnig.FkUleL5MpJo4BwpFU53FRiRxkf4t1fC',
        '超级管理员',
        '',
        'active',
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- 为管理员用户分配超级管理员角色
INSERT INTO
    user_roles (user_id, role_id, created_at)
VALUES (
        'user-admin',
        'role-admin',
        NOW()
    )
ON CONFLICT (user_id, role_id) DO NOTHING;

-- 打印提示信息
DO $$
BEGIN
    RAISE NOTICE '========================================';
    RAISE NOTICE '默认管理员账号已创建:';
    RAISE NOTICE '  用户名: admin';
    RAISE NOTICE '  密码: 123456';
    RAISE NOTICE '  邮箱: admin@ragljx.com';
    RAISE NOTICE '========================================';
    RAISE NOTICE '请在生产环境中立即修改默认密码！';
    RAISE NOTICE '========================================';
END $$;
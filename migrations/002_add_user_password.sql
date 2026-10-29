-- 首次初始化时在 001 之后执行；已有数据库需要手动执行此增量 SQL。
-- 旧用户的密码为空，不能登录；不设置共享默认密码。
BEGIN;
ALTER TABLE users ADD COLUMN password_hash TEXT;
COMMENT ON COLUMN users.password_hash IS 'bcrypt 密码哈希；旧用户为 NULL，不能登录';
COMMIT;

# 本地 PostgreSQL

使用 Docker Desktop 的 Linux 容器及 Docker Compose v2，在项目根目录执行以下命令。

## 配置与启动

项目根目录 `.env` 中配置数据库参数，参考 `.env.example`，保留已有的 LLM 配置：

```dotenv
POSTGRES_DB=finance
POSTGRES_USER=finance
POSTGRES_PASSWORD=填写本地数据库密码
POSTGRES_PORT=5432
```

```powershell
docker compose up -d postgres
```

首次启动会创建数据库，并执行 `migrations/001_create_conversation_tables.sql`，建立 users、sessions、messages 三张表。配置使用 PostgreSQL 18 官方镜像，数据保存在 `postgres_data` 命名卷中，挂载到 `/var/lib/postgresql`，实际数据目录默认为 `/var/lib/postgresql/18/docker`。容器和数据库默认时区均为 `Asia/Shanghai`。

如果之前已用 PostgreSQL 17 初始化过该数据卷，不能直接使用 18 读取旧数据，也不能仅靠调整目录升级。需要先用 17 导出数据，再导入使用独立新数据卷的 18 实例，或按官方流程使用 pg_upgrade。当前配置保留原有卷名，不自动迁移或删除数据。

初始化脚本只在数据目录为空时执行，重启不会重复建表。后续修改 SQL 文件不会自动更新已有表，需要另行执行增量迁移。初始化失败后应先排查日志，避免将已有数据目录误当成成功初始化。

## 连接参数

| 参数 | 值 |
| --- | --- |
| 主机 | `127.0.0.1` |
| 端口 | `.env` 的 POSTGRES_PORT，默认 `5432` |
| 数据库 | `.env` 的 POSTGRES_DB，默认 `finance` |
| 用户 | `.env` 的 POSTGRES_USER，默认 `finance` |
| 密码 | `.env` 的 POSTGRES_PASSWORD |
| 时区 | `Asia/Shanghai` |

端口仅监听本机。上述用户由官方镜像创建为数据库管理员，用于本地开发。当前 Go 服务尚未接入数据库，这些变量目前只供 Compose 使用。

用户名、数据库名和密码仅在首次初始化时生效；已有数据卷时修改 `.env` 不会修改库内账号或密码。

## 日常管理

```powershell
# 查看容器状态
docker compose ps

# 查看数据库日志
docker compose logs --tail 100 postgres

# 停止服务，保留数据卷
docker compose down
```

保留数据时不要给 down 加 `-v`，该选项会删除数据卷。

官方镜像及初始化行为参考：[PostgreSQL Docker 官方镜像说明](https://hub.docker.com/_/postgres)。

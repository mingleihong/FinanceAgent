# 智投顾问 · 前端

金融建议 Agent 的 React 注册、登录与聊天页面。

## 技术栈

- Vite + React 19 + TypeScript
- react-markdown + remark-gfm：渲染助手返回的 Markdown（列表 / 表格等）
- 开发模式下通过 Vite 代理将 `/api` 转发到 Go 后端，无跨域问题

## 运行

```bash
# 1. 先启动 Go 后端（默认 localhost:8080，见仓库根目录）
cd .. && go run .

# 2. 安装依赖并启动前端
npm install
npm run dev   # 打开 http://localhost:5173
```

## 构建

```bash
npm run build   # 产物输出到 dist/
```

## Docker Compose 部署

在项目根目录配置好现有 `.env` 后执行：

```bash
docker compose up -d --build
```

打开 `http://localhost`。前端镜像使用 Node 构建静态文件，由 Nginx 提供页面，
并将 `/api/` 请求代理到 `backend:8080`，保留聊天的 SSE 流式输出。
只有前端向宿主机开放 80 端口，后端和数据库通过 Compose 内部网络访问。
后端继续使用现有 `.env` 和 `config.yaml`，无需新增环境变量。
若修改 `config.yaml` 中的 `app.port`，需同步修改 `nginx.conf` 的后端端口并重新构建。

## 注册与登录

- `POST /api/register` 和 `POST /api/login` 均发送 `{ "username": "...", "password": "..." }`。
- 使用后端统一响应格式 `{ "code": 0, "message": "success", "data": ... }`；注册成功后切换到登录页，保留用户名并清除密码。
- 用户名去除首尾空白后须为 1–64 个 Unicode 字符，不能包含空字符；密码必填且最多 72 个 UTF-8 字节，与后端约束一致。
- 登录读取 `data.access_token` 与 `data.token_type`，在当前浏览器的 `localStorage` 中以 `finance.session` 为键保存凭证；刷新、关闭标签页或重启浏览器后可恢复登录，直到 Token 过期或主动退出，不保存密码。禁用浏览器存储时仅保留内存登录态。
- 前端依据 JWT 的 `exp` 处理过期，聊天请求携带 `Authorization: Bearer <token>`；收到 HTTP 401 时返回登录页。退出或过期时终止当前聊天请求并清空当前对话。
- 后端 `/api/chat` 通过 JWT 中间件校验 Bearer Token，限定 HS256 算法，校验签名、必填的过期时间与用户身份字段；缺少凭证、签名无效或过期时返回统一 JSON 错误（HTTP 401），不会启动 SSE。注册和登录接口保持公开。退出登录清除本地凭证，后端目前没有 Token 撤销接口。

## 聊天说明

后端 `POST /api/chat` 接收 `{ "prompt": "..." }`，默认通过 SSE 流式
返回 `reasoning_content`（思考内容）、`content`（回答增量）、`done` 或 `error` 事件。
事件名携带类型，JSON 中没有 `type` 字段：

| SSE 事件 | JSON 数据 | 前端行为 |
| --- | --- | --- |
| `reasoning_content` | `{ "content": "思考增量" }` | 追加到思考区 |
| `content` | `{ "content": "回答增量" }` | 追加到回答区 |
| `done` | `{}` | 结束生成状态，允许继续发送 |
| `error` | `{ "message": "错误说明" }` | 显示错误，结束生成状态 |

聊天请求入口为 `streamMessage`，回调参数为 `(event, content)`。
请求参数校验失败时读取 HTTP 错误响应的 `message`；流中断且未收到 `done`
时显示回答可能不完整，并保留已经收到的思考和正文。

思考内容仅来自模型供应商返回的 `reasoning_content` 字符串字段；
模型及接口必须支持并启用思考输出。未返回这些字段时仅显示流式回答。
工具调用结果会继续交给模型生成回答。当前仍无服务端对话记忆，
前端在浏览器内维护完整会话展示。

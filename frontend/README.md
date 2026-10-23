# 智投顾问 · 前端

金融建议 Agent 的 React 聊天页面。

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

## 说明

后端 `POST /api/chat` 接收 `{ "prompt": "..." }`，默认通过 SSE 流式
返回 `reasoning_content`（思考内容）、`content`（回答增量）、`done` 或 `error` 事件。
事件名携带类型，JSON 中没有 `type` 字段：

| SSE 事件 | JSON 数据 | 前端行为 |
| --- | --- | --- |
| `reasoning_content` | `{ "content": "思考增量" }` | 追加到思考区 |
| `content` | `{ "content": "回答增量" }` | 追加到回答区 |
| `done` | `{}` | 结束生成状态，允许继续发送 |
| `error` | `{ "message": "错误说明" }` | 显示错误，结束生成状态 |

前端只有 `streamMessage` 一个请求入口，回调参数为 `(event, content)`。
请求参数校验失败时读取 HTTP 错误响应的 `message`；流中断且未收到 `done`
时显示回答可能不完整，并保留已经收到的思考和正文。

思考内容仅来自模型供应商返回的 `reasoning_content` 字符串字段；
模型及接口必须支持并启用思考输出。未返回这些字段时仅显示流式回答。
工具调用结果会继续交给模型生成回答。当前仍无服务端对话记忆，
前端在浏览器内维护完整会话展示。

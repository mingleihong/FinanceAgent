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

当前后端 `POST /api/chat` 为一次性问答（无对话记忆、非流式），
前端在浏览器内维护完整会话展示。后端后续若支持流式 / 带上下文，
只需扩展 `src/api/chat.ts`。

import type { ChatResponse } from '../types'

const API_ENDPOINT = '/api/chat'

export class ChatError extends Error {
  constructor(message: string) {
    super(message)
    this.name = 'ChatError'
  }
}

/**
 * 调用后端聊天接口。
 * 返回助手最终回答文本；失败时抛出 ChatError。
 */
export async function sendMessage(
  prompt: string,
  signal?: AbortSignal,
): Promise<string> {
  let res: Response
  try {
    res = await fetch(API_ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ prompt }),
      signal,
    })
  } catch (err) {
    if (err instanceof DOMException && err.name === 'AbortError') {
      throw err
    }
    throw new ChatError('无法连接到服务器，请确认后端服务已启动（localhost:8080）')
  }

  let payload: unknown
  try {
    payload = await res.json()
  } catch {
    throw new ChatError('服务器返回了无法解析的数据')
  }

  const body = payload as ChatResponse
  if (!res.ok) {
    throw new ChatError(body.message ?? `请求失败（HTTP ${res.status}）`)
  }
  if (body.code !== 0) {
    throw new ChatError(body.message ?? '服务器返回了错误响应')
  }

  return body.data ?? ''
}

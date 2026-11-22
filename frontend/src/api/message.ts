import { apiErrorMessage, AuthenticationError } from './error'

/** 后端消息 DTO，ID 为字符串以避免大整数精度丢失。 */
export interface HistoryMessage {
  id: string
  session_id: string
  role: string
  content: string
  reasoning_content?: string | null
  created_at: string
}

interface Envelope<T> {
  code?: unknown
  message?: unknown
  data?: T
}

function parseData<T>(body: unknown, fallback: string): T {
  const envelope = body as Envelope<T> | null
  if (envelope?.code !== 0) {
    const message = apiErrorMessage(body, fallback)
    if (envelope?.code === 10002) throw new AuthenticationError(message)
    throw new Error(message)
  }
  if (envelope.data === undefined) throw new Error(fallback)
  return envelope.data
}

/** 读取指定会话的历史消息；会话归属由后端校验。 */
export async function listMessages(
  sessionId: string,
  accessToken: string,
  signal?: AbortSignal,
): Promise<HistoryMessage[]> {
  let res: Response
  try {
    res = await fetch(`/api/sessions/${encodeURIComponent(sessionId)}/messages`, {
      headers: { Authorization: `Bearer ${accessToken}` },
      signal,
      cache: 'no-store',
    })
  } catch (err) {
    if (err instanceof DOMException && err.name === 'AbortError') throw err
    throw new Error('无法连接到服务器，请确认后端服务已启动')
  }

  const body = await res.json().catch(() => null)
  if (res.status === 401) {
    throw new AuthenticationError(apiErrorMessage(body, '登录已失效，请重新登录。'))
  }
  const data = parseData<HistoryMessage[]>(body, '加载聊天记录失败，请稍后重试。')
  if (!Array.isArray(data)) throw new Error('聊天记录返回格式异常')
  return data
}

import { apiErrorMessage, AuthenticationError } from './error'

export interface SessionItem {
  id: string
  title: string
  created_at: string
  updated_at: string
}

export interface SessionPage {
  list: SessionItem[]
  total: number
  page: number
  page_size: number
}

interface Envelope<T> {
  code?: unknown
  message?: unknown
  data?: T
}

async function parse<T>(res: Response, fallback: string): Promise<T> {
  const body = await res.json().catch(() => null)
  if (!res.ok || (body as Envelope<T> | null)?.code !== 0) {
    const message = apiErrorMessage(body, res.status === 401 ? '登录已失效，请重新登录。' : fallback)
    if (res.status === 401 || (body as Envelope<T> | null)?.code === 10002) {
      throw new AuthenticationError(message)
    }
    throw new Error(message)
  }
  const data = (body as Envelope<T>).data
  if (data === undefined) throw new Error(fallback)
  return data
}

export async function listSessions(
  page: number,
  pageSize: number,
  accessToken: string,
  signal?: AbortSignal,
): Promise<SessionPage> {
  let res: Response
  try {
    res = await fetch(`/api/sessions?page=${page}&page_size=${pageSize}`, {
      headers: { Authorization: `Bearer ${accessToken}` },
      signal,
      cache: 'no-store',
    })
  } catch (err) {
    if (err instanceof DOMException && err.name === 'AbortError') throw err
    throw new Error('无法连接到服务器，请确认后端服务已启动')
  }
  const data = await parse<SessionPage>(res, '加载会话列表失败，请稍后重试。')
  if (!data || !Array.isArray(data.list)) throw new Error('会话列表返回格式异常')
  return data
}

export async function createSession(accessToken: string): Promise<SessionItem> {
  let res: Response
  try {
    res = await fetch('/api/sessions', {
      method: 'POST',
      headers: { Authorization: `Bearer ${accessToken}` },
      cache: 'no-store',
    })
  } catch {
    throw new Error('无法连接到服务器，请确认后端服务已启动')
  }
  const data = await parse<SessionItem>(res, '创建会话失败，请稍后重试。')
  if (!data || typeof data.id !== 'string') throw new Error('创建会话返回格式异常')
  return data
}

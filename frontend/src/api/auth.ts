export interface Session {
  username: string
  accessToken: string
  expiresAt: number
}

const SESSION_KEY = 'finance.session'

const errorMessages: Record<string, string> = {
  'invalid login request': '请填写用户名和密码。',
  'invalid registration request': '请填写完整的注册信息。',
  'invalid username or password': '用户名或密码不正确。',
  'username already exists': '该用户名已被注册，请更换用户名或直接登录。',
  'login failed': '登录失败，请稍后重试。',
  'registration failed': '注册失败，请稍后重试。',
  'invalid registration: username must contain 1 to 64 characters and no null characters': '用户名须为 1–64 个字符，且不能包含空字符。',
  'invalid registration: password must not exceed 72 bytes': '密码不能超过 72 字节。',
}

async function request(path: 'login' | 'register', username: string, password: string, signal: AbortSignal): Promise<unknown> {
  let response: Response
  try {
    response = await fetch(`/api/${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
      signal,
      cache: 'no-store',
    })
  } catch (error) {
    if (signal.aborted) throw error
    throw new Error('无法连接服务器，请检查网络后重试。')
  }
  const body = await response.json().catch(() => null)
  if (!response.ok || body?.code !== 0) {
    const message = typeof body?.message === 'string' ? errorMessages[body.message] : undefined
    throw new Error(message ?? (response.status === 429 ? '操作过于频繁，请稍后重试。' : '请求失败，请稍后重试。'))
  }
  return body.data
}

export async function register(username: string, password: string, signal: AbortSignal): Promise<void> {
  await request('register', username, password, signal)
}

export async function login(username: string, password: string, signal: AbortSignal): Promise<Session> {
  const data = await request('login', username, password, signal) as { access_token?: unknown; token_type?: unknown } | null
  if (typeof data?.access_token !== 'string' || data.token_type !== 'Bearer') {
    throw new Error('服务器返回的登录信息不完整，请重试。')
  }
  // Decode expiry for UI lifecycle only; token verification belongs to the server.
  try {
    const payload = data.access_token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
    const claims = JSON.parse(atob(payload.padEnd(Math.ceil(payload.length / 4) * 4, '=')))
    const expiresAt = claims.exp * 1000
    if (typeof claims.exp !== 'number' || !Number.isFinite(expiresAt) || expiresAt <= Date.now()) throw new Error()
    return { username, accessToken: data.access_token, expiresAt }
  } catch {
    throw new Error('登录凭证无效或已过期，请重新登录。')
  }
}

export function readSession(): Session | null {
  try {
    const session = JSON.parse(localStorage.getItem(SESSION_KEY) ?? 'null')
    if (typeof session?.username === 'string' && typeof session.accessToken === 'string' && session.accessToken &&
        typeof session.expiresAt === 'number' && session.expiresAt > Date.now()) return session
    localStorage.removeItem(SESSION_KEY)
  } catch { /* Storage may be unavailable in private browsing. */ }
  return null
}

export function saveSession(session: Session | null) {
  try {
    if (session) localStorage.setItem(SESSION_KEY, JSON.stringify(session))
    else localStorage.removeItem(SESSION_KEY)
  } catch { /* Keep the current in-memory session if browser storage is disabled. */ }
}

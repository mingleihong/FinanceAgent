const errorMessages: Record<number, string> = {
  10000: '服务器内部错误，请稍后重试。',
  10002: '登录已失效，请重新登录。',
  10003: '无权限访问。',
  10004: '资源不存在。',
  10005: '请求方法不支持。',
  20001: '用户不存在。',
  20002: '该用户名已被注册，请更换用户名或直接登录。',
  20003: '用户名或密码不正确。',
  30001: '回答生成失败，请稍后重试。',
}

// Validation errors can carry a more specific, public message from the server.
export function apiErrorMessage(body: unknown, fallback: string): string {
  if (typeof body !== 'object' || body === null) return fallback

  const { code, message } = body as { code?: unknown; message?: unknown }
  if (typeof code === 'number' && errorMessages[code]) return errorMessages[code]
  return typeof message === 'string' && message.trim() ? message : fallback
}

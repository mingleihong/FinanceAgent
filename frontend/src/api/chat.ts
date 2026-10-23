import type { ChatEvent } from '../types'

const API_ENDPOINT = '/api/chat'

/** Reads SSE frames independently of network chunk and UTF-8 boundaries. */
export async function streamMessage(
  prompt: string,
  onDelta: (event: ChatEvent, content: string) => void,
  signal?: AbortSignal,
): Promise<void> {
  let res: Response
  try {
    res = await fetch(API_ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Accept: 'text/event-stream' },
      body: JSON.stringify({ prompt }),
      signal,
    })
  } catch (err) {
    if (err instanceof DOMException && err.name === 'AbortError') throw err
    throw new Error('无法连接到服务器，请确认后端服务已启动')
  }

  if (!res.ok) {
    const body = await res.json().catch(() => null)
    throw new Error(typeof body?.message === 'string' ? body.message : `请求失败（HTTP ${res.status}）`)
  }
  if (!res.body || !res.headers.get('content-type')?.includes('text/event-stream')) {
    throw new Error('服务器未返回流式响应')
  }

  const reader = res.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  let completed = false

  function consume(frame: string) {
    let event = 'message'
    const data: string[] = []
    for (const line of frame.split('\n')) {
      if (line.startsWith('event:')) event = line.slice(6).trim()
      if (line.startsWith('data:')) data.push(line.slice(5).replace(/^ /, ''))
    }
    if (data.length === 0) return
    if (event === 'done') {
      completed = true
      return
    }
    if (event !== 'content' && event !== 'reasoning_content' && event !== 'error') return
    let payload: unknown
    try {
      payload = JSON.parse(data.join('\n'))
    } catch {
      throw new Error('服务器返回了无法解析的数据')
    }
    const body = payload as { content?: unknown; message?: unknown } | null
    if (event === 'error') {
      throw new Error(typeof body?.message === 'string' ? body.message : '回答生成失败，请稍后重试')
    }
    if (typeof body?.content !== 'string') {
      throw new Error('服务器返回了无效的流式内容')
    }
    onDelta(event, body.content)
  }

  try {
    while (!completed) {
      const { value, done } = await reader.read()
      buffer += decoder.decode(value, { stream: !done })
      buffer = buffer.replace(/\r\n/g, '\n')
      let boundary: number
      while (!completed && (boundary = buffer.indexOf('\n\n')) !== -1) {
        consume(buffer.slice(0, boundary))
        buffer = buffer.slice(boundary + 2)
      }
      if (done) {
        if (!completed) throw new Error('连接已中断，回答可能不完整，请重试')
        break
      }
    }
  } finally {
    await reader.cancel().catch(() => {})
    reader.releaseLock()
  }
}

import { useCallback, useEffect, useRef, useState } from 'react'
import { AuthenticationError, streamMessage } from './api/chat'
import type { ChatTurn } from './api/chat'
import { createSession, listSessions } from './api/session'
import type { SessionPage } from './api/session'
import { listMessages } from './api/message'
import type { Message } from './types'
import Header from './components/Header'
import Welcome from './components/Welcome'
import MessageList from './components/MessageList'
import ChatInput from './components/ChatInput'
import SessionList from './components/SessionList'
import AuthPage from './components/AuthPage'
import { readSession, saveSession } from './api/auth'
import type { Session } from './api/auth'

const PAGE_SIZE = 10

export default function App() {
  const [session, setSession] = useState<Session | null>(readSession)
  const [authNotice, setAuthNotice] = useState('')
  const [messages, setMessages] = useState<Message[]>([])
  // null 表示尚未绑定到服务端的「新聊天」页面，不产生任何会话记录。
  const [activeSessionId, setActiveSessionId] = useState<string | null>(null)
  const [historyLoading, setHistoryLoading] = useState(false)
  const [historyError, setHistoryError] = useState('')
  const scrollRef = useRef<HTMLDivElement>(null)
  const requestRef = useRef<AbortController | null>(null)
  const historyRef = useRef<AbortController | null>(null)
  const pending = messages[messages.length - 1]?.streaming ?? false

  const [sessionPage, setSessionPage] = useState<SessionPage | null>(null)
  const [sessionListLoading, setSessionListLoading] = useState(false)
  const [sessionListError, setSessionListError] = useState('')

  useEffect(() => () => {
    requestRef.current?.abort()
    historyRef.current?.abort()
  }, [])

  function logout(notice = '') {
    requestRef.current?.abort()
    historyRef.current?.abort()
    requestRef.current = null
    historyRef.current = null
    saveSession(null)
    setSession(null)
    setMessages([])
    setActiveSessionId(null)
    setHistoryLoading(false)
    setHistoryError('')
    setSessionPage(null)
    setSessionListError('')
    setAuthNotice(notice)
  }

  // 会话列表：仅登录态下加载；401 统一登出，其余错误就地展示。
  const loadSessionPage = useCallback(async (page: number, token: string) => {
    setSessionListLoading(true)
    setSessionListError('')
    try {
      const data = await listSessions(page, PAGE_SIZE, token)
      setSessionPage(data)
    } catch (err) {
      if (err instanceof AuthenticationError) {
        logout(err.message)
        return
      }
      setSessionListError(err instanceof Error ? err.message : '加载会话列表失败')
    } finally {
      setSessionListLoading(false)
    }
  }, [])

  useEffect(() => {
    if (session) void loadSessionPage(1, session.accessToken)
  }, [session, loadSessionPage])

  async function changeSessionPage(page: number) {
    if (!session || !sessionPage || page < 1 || page > Math.ceil(sessionPage.total / PAGE_SIZE)) return
    await loadSessionPage(page, session.accessToken)
  }

  // 点会话：切换视图并读取该会话的历史消息，不产生新的会话记录。
  async function selectSession(id: string) {
    if (!session) return
    if (id === activeSessionId && !historyError) return
    requestRef.current?.abort()
    requestRef.current = null
    historyRef.current?.abort()
    const controller = new AbortController()
    historyRef.current = controller
    setActiveSessionId(id)
    setMessages([])
    setHistoryError('')
    setHistoryLoading(true)
    try {
      const history = await listMessages(id, session.accessToken, controller.signal)
      if (historyRef.current !== controller) return
      setMessages(history
        .filter((item) => item.role === 'user' || item.role === 'assistant')
        .map((item) => ({
          id: item.id,
          role: item.role as 'user' | 'assistant',
          content: item.content,
          ...(item.reasoning_content ? { reasoning_content: item.reasoning_content } : {}),
        })))
    } catch (err) {
      if (controller.signal.aborted) return
      if (err instanceof AuthenticationError) {
        logout(err.message)
        return
      }
      setHistoryError(err instanceof Error ? err.message : '加载聊天记录失败')
    } finally {
      if (historyRef.current === controller) {
        historyRef.current = null
        setHistoryLoading(false)
      }
    }
  }

  // 点新会话：只进入空白聊天页，真正的会话在首次发送消息时才创建。
  function startNewChat() {
    requestRef.current?.abort()
    historyRef.current?.abort()
    requestRef.current = null
    historyRef.current = null
    setActiveSessionId(null)
    setMessages([])
    setHistoryLoading(false)
    setHistoryError('')
  }

  useEffect(() => {
    if (!session) return
    // Recheck on focus as background tabs can delay timers.
    function checkExpiry() {
      if (session && session.expiresAt <= Date.now()) logout('登录已过期，请重新登录。')
    }
    const timer = window.setInterval(checkExpiry, 1000)
    window.addEventListener('focus', checkExpiry)
    return () => {
      window.clearInterval(timer)
      window.removeEventListener('focus', checkExpiry)
    }
  }, [session])

  // 新消息或思考中状态变化时自动滚动到底部。
  useEffect(() => {
    const el = scrollRef.current
    if (el) el.scrollTop = el.scrollHeight
  }, [messages, historyLoading])

  function toTurns(list: Message[]): ChatTurn[] {
    return list
      .filter((item) => (item.role === 'user' || item.role === 'assistant') && !item.streaming && !item.error && item.content.trim())
      .map((item) => ({
        role: item.role as 'user' | 'assistant',
        content: item.content,
        ...(item.reasoning_content ? { reasoning_content: item.reasoning_content } : {}),
      }))
  }

  async function send(prompt: string) {
    prompt = prompt.trim()
    if (!prompt || requestRef.current || !session || historyLoading) return
    if (session.expiresAt <= Date.now()) {
      logout('登录已过期，请重新登录。')
      return
    }
    const controller = new AbortController()
    requestRef.current = controller
    const userMessage: Message = {
      id: crypto.randomUUID(),
      role: 'user',
      content: prompt,
    }
    const assistantID = crypto.randomUUID()
    setMessages((prev) => [...prev, userMessage, {
      id: assistantID,
      role: 'assistant',
      content: '',
      reasoning_content: '',
      streaming: true,
    }])

    let targetSessionId = activeSessionId
    try {
      // 新聊天页首次发送：先建会话拿到 session_id，再发送对话内容。
      if (!targetSessionId) {
        const created = await createSession(session.accessToken)
        if (controller.signal.aborted) return
        targetSessionId = created.id
        setActiveSessionId(created.id)
        void loadSessionPage(1, session.accessToken)
      }

      await streamMessage({
        prompt,
        sessionId: targetSessionId,
        messages: toTurns([...messages, userMessage]),
      }, (event, content) => {
        setMessages((prev) => prev.map((message) => {
          if (message.id !== assistantID) return message
          return { ...message, [event]: (message[event] ?? '') + content }
        }))
      }, controller.signal, session.accessToken)

      // 消息已由后端落库，刷新列表中的更新时间与排序。
      void loadSessionPage(sessionPage?.page ?? 1, session.accessToken)
    } catch (err) {
      if (controller.signal.aborted) return
      if (err instanceof AuthenticationError) {
        logout(err.message)
        return
      }
      const text = err instanceof Error ? err.message : '出错了，请稍后重试'
      setMessages((prev) => prev.map((message) => message.id === assistantID
        ? { ...message, error: text }
        : message))
    } finally {
      setMessages((prev) => prev.map((message) => message.id === assistantID
        ? { ...message, streaming: false }
        : message))
      if (requestRef.current === controller) requestRef.current = null
    }
  }

  const showWelcome = !historyLoading && !historyError && messages.length === 0

  if (!session) {
    return <AuthPage notice={authNotice} onLogin={(nextSession) => {
      saveSession(nextSession)
      setSession(nextSession)
      setAuthNotice('')
    }} />
  }

  return (
    <div className="app">
      <Header username={session.username} onLogout={() => logout()} />
      <div className="workspace">
        <aside className="sidebar">
          <SessionList
            sessions={sessionPage?.list ?? []}
            total={sessionPage?.total ?? 0}
            page={sessionPage?.page ?? 1}
            pageSize={PAGE_SIZE}
            loading={sessionListLoading}
            error={sessionListError}
            activeId={activeSessionId}
            onPageChange={(p) => void changeSessionPage(p)}
            onRefresh={() => void loadSessionPage(sessionPage?.page ?? 1, session.accessToken)}
            onSelect={(id) => void selectSession(id)}
            onNewChat={startNewChat}
          />
        </aside>
        <main className="chat-main">
          <div className="chat" ref={scrollRef}>
            {historyLoading ? (
              <p className="chat-placeholder">正在加载聊天记录…</p>
            ) : historyError ? (
              <div className="chat-placeholder">
                <p className="chat-error-text">{historyError}</p>
                {activeSessionId && (
                  <button type="button" className="chat-retry" onClick={() => void selectSession(activeSessionId)}>
                    重试
                  </button>
                )}
              </div>
            ) : showWelcome ? (
              <Welcome onPick={send} />
            ) : (
              <MessageList messages={messages} />
            )}
          </div>
          <ChatInput disabled={pending || historyLoading} onSend={send} />
        </main>
      </div>
    </div>
  )
}

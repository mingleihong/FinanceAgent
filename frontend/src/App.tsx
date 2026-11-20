import { useCallback, useEffect, useRef, useState } from 'react'
import { AuthenticationError, streamMessage } from './api/chat'
import { createSession, listSessions } from './api/session'
import type { SessionPage } from './api/session'
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
  const scrollRef = useRef<HTMLDivElement>(null)
  const requestRef = useRef<AbortController | null>(null)
  const pending = messages[messages.length - 1]?.streaming ?? false

  const [sessionPage, setSessionPage] = useState<SessionPage | null>(null)
  const [sessionListLoading, setSessionListLoading] = useState(false)
  const [sessionListError, setSessionListError] = useState('')
  const [creatingSession, setCreatingSession] = useState(false)

  useEffect(() => () => requestRef.current?.abort(), [])

  function logout(notice = '') {
    requestRef.current?.abort()
    requestRef.current = null
    saveSession(null)
    setSession(null)
    setMessages([])
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

  async function createNewSession() {
    if (!session || creatingSession) return
    setCreatingSession(true)
    setSessionListError('')
    try {
      await createSession(session.accessToken)
      await loadSessionPage(1, session.accessToken)
    } catch (err) {
      if (err instanceof AuthenticationError) {
        logout(err.message)
        return
      }
      setSessionListError(err instanceof Error ? err.message : '创建会话失败')
    } finally {
      setCreatingSession(false)
    }
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
  }, [messages])

  async function send(prompt: string) {
    prompt = prompt.trim()
    if (!prompt || requestRef.current || !session) return
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

    try {
      await streamMessage(prompt, (event, content) => {
        setMessages((prev) => prev.map((message) => {
          if (message.id !== assistantID) return message
          return { ...message, [event]: (message[event] ?? '') + content }
        }))
      }, controller.signal, session.accessToken)
    } catch (err) {
      if (controller.signal.aborted) return
      if (err instanceof AuthenticationError) {
        logout(err.message)
        return
      }
      const text =
        err instanceof Error
          ? err.message
          : '出错了，请稍后重试'
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

  const showWelcome = messages.length === 0

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
            creating={creatingSession}
            onPageChange={(p) => void changeSessionPage(p)}
            onRefresh={() => void loadSessionPage(sessionPage?.page ?? 1, session.accessToken)}
            onCreate={() => void createNewSession()}
          />
        </aside>
        <main className="chat-main">
          <div className="chat" ref={scrollRef}>
            {showWelcome ? (
              <Welcome onPick={send} />
            ) : (
              <MessageList messages={messages} />
            )}
          </div>
          <ChatInput disabled={pending} onSend={send} />
        </main>
      </div>
    </div>
  )
}

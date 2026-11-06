import { useEffect, useRef, useState } from 'react'
import { AuthenticationError, streamMessage } from './api/chat'
import type { Message } from './types'
import Header from './components/Header'
import Welcome from './components/Welcome'
import MessageList from './components/MessageList'
import ChatInput from './components/ChatInput'
import AuthPage from './components/AuthPage'
import { readSession, saveSession } from './api/auth'
import type { Session } from './api/auth'

export default function App() {
  const [session, setSession] = useState<Session | null>(readSession)
  const [authNotice, setAuthNotice] = useState('')
  const [messages, setMessages] = useState<Message[]>([])
  const scrollRef = useRef<HTMLDivElement>(null)
  const requestRef = useRef<AbortController | null>(null)
  const pending = messages[messages.length - 1]?.streaming ?? false

  useEffect(() => () => requestRef.current?.abort(), [])

  function logout(notice = '') {
    requestRef.current?.abort()
    requestRef.current = null
    saveSession(null)
    setSession(null)
    setMessages([])
    setAuthNotice(notice)
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
      <div className="chat" ref={scrollRef}>
        {showWelcome ? (
          <Welcome onPick={send} />
        ) : (
          <MessageList messages={messages} />
        )}
      </div>
      <ChatInput disabled={pending} onSend={send} />
    </div>
  )
}

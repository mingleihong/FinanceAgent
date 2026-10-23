import { useEffect, useRef, useState } from 'react'
import { streamMessage } from './api/chat'
import type { Message } from './types'
import Header from './components/Header'
import Welcome from './components/Welcome'
import MessageList from './components/MessageList'
import ChatInput from './components/ChatInput'

export default function App() {
  const [messages, setMessages] = useState<Message[]>([])
  const scrollRef = useRef<HTMLDivElement>(null)
  const requestRef = useRef<AbortController | null>(null)
  const pending = messages[messages.length - 1]?.streaming ?? false

  useEffect(() => () => requestRef.current?.abort(), [])

  // 新消息或思考中状态变化时自动滚动到底部。
  useEffect(() => {
    const el = scrollRef.current
    if (el) el.scrollTop = el.scrollHeight
  }, [messages])

  async function send(prompt: string) {
    prompt = prompt.trim()
    if (!prompt || requestRef.current) return
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
      }, controller.signal)
    } catch (err) {
      if (controller.signal.aborted) return
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
      requestRef.current = null
    }
  }

  const showWelcome = messages.length === 0

  return (
    <div className="app">
      <Header />
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

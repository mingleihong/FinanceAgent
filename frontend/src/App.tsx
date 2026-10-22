import { useEffect, useRef, useState } from 'react'
import { ChatError, sendMessage } from './api/chat'
import type { Message } from './types'
import Header from './components/Header'
import Welcome from './components/Welcome'
import MessageList from './components/MessageList'
import ChatInput from './components/ChatInput'

export default function App() {
  const [messages, setMessages] = useState<Message[]>([])
  const [pending, setPending] = useState(false)
  const scrollRef = useRef<HTMLDivElement>(null)

  // 新消息或思考中状态变化时自动滚动到底部。
  useEffect(() => {
    const el = scrollRef.current
    if (el) el.scrollTop = el.scrollHeight
  }, [messages, pending])

  async function send(prompt: string) {
    const userMessage: Message = {
      id: crypto.randomUUID(),
      role: 'user',
      content: prompt,
    }
    setMessages((prev) => [...prev, userMessage])
    setPending(true)

    try {
      const reply = await sendMessage(prompt)
      setMessages((prev) => [
        ...prev,
        {
          id: crypto.randomUUID(),
          role: 'assistant',
          content: reply,
        },
      ])
    } catch (err) {
      const text =
        err instanceof ChatError
          ? err.message
          : '出错了，请稍后重试'
      setMessages((prev) => [
        ...prev,
        {
          id: crypto.randomUUID(),
          role: 'assistant',
          content: text,
          error: true,
        },
      ])
    } finally {
      setPending(false)
    }
  }

  const showWelcome = messages.length === 0 && !pending

  return (
    <div className="app">
      <Header />
      <div className="chat" ref={scrollRef}>
        {showWelcome ? (
          <Welcome onPick={send} />
        ) : (
          <MessageList messages={messages} pending={pending} />
        )}
      </div>
      <ChatInput disabled={pending} onSend={send} />
    </div>
  )
}

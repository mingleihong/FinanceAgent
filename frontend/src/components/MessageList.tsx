import type { Message } from '../types'
import MessageItem from './MessageItem'
import Thinking from './Thinking'

interface MessageListProps {
  messages: Message[]
  pending: boolean
}

export default function MessageList({ messages, pending }: MessageListProps) {
  return (
    <div className="message-list">
      {messages.map((message) => (
        <MessageItem key={message.id} message={message} />
      ))}
      {pending && <Thinking />}
    </div>
  )
}

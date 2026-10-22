import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import type { Message } from '../types'

interface MessageItemProps {
  message: Message
}

export default function MessageItem({ message }: MessageItemProps) {
  if (message.role === 'user') {
    return (
      <div className="message user">
        <div className="bubble user-bubble">{message.content}</div>
      </div>
    )
  }

  return (
    <div className={`message assistant${message.error ? ' error' : ''}`}>
      <span className="avatar" aria-hidden="true">
        🤖
      </span>
      <div className="bubble">
        <ReactMarkdown
          remarkPlugins={[remarkGfm]}
          components={{
            a: (props) => (
              <a {...props} target="_blank" rel="noreferrer" />
            ),
          }}
        >
          {message.content}
        </ReactMarkdown>
        {message.error && (
          <p className="error-note">（请求失败，请稍后重试）</p>
        )}
      </div>
    </div>
  )
}

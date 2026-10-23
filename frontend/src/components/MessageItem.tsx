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
        {message.reasoning_content && (
          <details className="reasoning" open>
            <summary>思考内容</summary>
            <div className="reasoning-content">{message.reasoning_content}</div>
          </details>
        )}
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
        {message.streaming && (
          <p className="stream-status" role="status">
            {message.content ? '正在生成回答…' : message.reasoning_content ? '思考中…' : '正在等待模型响应…'}
          </p>
        )}
        {message.error && (
          <p className="error-note" role="alert">{message.error}</p>
        )}
      </div>
    </div>
  )
}

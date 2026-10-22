import { useRef, useState, type KeyboardEvent } from 'react'

interface ChatInputProps {
  disabled: boolean
  onSend: (text: string) => void
}

export default function ChatInput({ disabled, onSend }: ChatInputProps) {
  const [value, setValue] = useState('')
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  function submit() {
    const text = value.trim()
    if (!text || disabled) return
    onSend(text)
    setValue('')
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto'
    }
  }

  function handleInput() {
    const el = textareaRef.current
    if (!el) return
    el.style.height = 'auto'
    el.style.height = `${Math.min(el.scrollHeight, 160)}px`
  }

  function handleKeyDown(e: KeyboardEvent<HTMLTextAreaElement>) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      submit()
    }
  }

  return (
    <div className="input-bar">
      <div className="input-box">
        <textarea
          ref={textareaRef}
          className="input-textarea"
          placeholder={disabled ? '助手正在思考中…' : '输入你的理财问题，Enter 发送，Shift+Enter 换行'}
          value={value}
          disabled={disabled}
          rows={1}
          onChange={(e) => setValue(e.target.value)}
          onInput={handleInput}
          onKeyDown={handleKeyDown}
        />
        <button
          type="button"
          className="send-button"
          onClick={submit}
          disabled={disabled || !value.trim()}
          aria-label="发送"
        >
          <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
            <path
              fill="currentColor"
              d="M3.4 20.4l17.4-7.5c.8-.34.8-1.6 0-1.94L3.4 3.4c-.66-.3-1.4.24-1.4.96v5.7c0 .5.36.94.85 1l10.3 1.4-10.3 1.4c-.5.08-.85.5-.85 1v5.7c0 .72.74 1.26 1.4.96z"
            />
          </svg>
        </button>
      </div>
      <p className="input-hint">内容由 AI 生成，仅供参考，不构成投资建议</p>
    </div>
  )
}

import { useEffect, useRef, useState } from 'react'
import type { FormEvent } from 'react'
import { login, register } from '../api/auth'
import type { Session } from '../api/auth'

interface Props {
  onLogin: (session: Session) => void
  notice: string
}

export default function AuthPage({ onLogin, notice }: Props) {
  const [mode, setMode] = useState<'login' | 'register'>('login')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirmation, setConfirmation] = useState('')
  const [visible, setVisible] = useState(false)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [message, setMessage] = useState(notice)
  const requestRef = useRef<AbortController | null>(null)
  const passwordRef = useRef<HTMLInputElement>(null)
  const isRegister = mode === 'register'

  useEffect(() => () => requestRef.current?.abort(), [])

  function switchMode() {
    setMode(isRegister ? 'login' : 'register')
    setPassword('')
    setConfirmation('')
    setVisible(false)
    setError('')
    setMessage('')
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (requestRef.current) return
    setError('')
    setMessage('')
    const name = username.trim()
    if (!name || Array.from(name).length > 64 || name.includes('\0')) {
      setError('用户名须为 1–64 个字符，且不能包含空字符。')
      return
    }
    if (!password || new TextEncoder().encode(password).length > 72) {
      setError('请输入密码，长度不能超过 72 字节（中文通常占 3 字节）。')
      return
    }
    if (isRegister && password !== confirmation) {
      setError('两次输入的密码不一致。')
      return
    }
    const controller = new AbortController()
    requestRef.current = controller
    setBusy(true)
    try {
      if (isRegister) {
        await register(name, password, controller.signal)
        setUsername(name)
        setPassword('')
        setConfirmation('')
        setVisible(false)
        setMode('login')
        setMessage('注册成功，请使用刚刚创建的账号登录。')
        passwordRef.current?.focus()
      } else {
        onLogin(await login(name, password, controller.signal))
      }
    } catch (err) {
      if (!controller.signal.aborted) setError(err instanceof Error ? err.message : '操作失败，请稍后重试。')
    } finally {
      requestRef.current = null
      if (!controller.signal.aborted) setBusy(false)
    }
  }

  return (
    <main className="auth-page">
      <section className="auth-intro" aria-label="智投顾问介绍">
        <div className="auth-brand"><span className="brand-mark" aria-hidden="true">💹</span>智投顾问</div>
        <div className="auth-intro-content">
          <span className="auth-eyebrow">YOUR FINANCIAL COMPANION</span>
          <h1>让每一次思考，<br />更接近明智决策。</h1>
          <p>从理解金融概念到梳理投资思路，<br />与 AI 顾问一起，让复杂问题变得清晰。</p>
          <div className="auth-topics"><span>金融知识</span><span>投资思路</span><span>风险意识</span></div>
        </div>
        <p className="auth-disclaimer">AI 内容仅供参考，不构成投资建议。</p>
      </section>
      <section className="auth-panel" aria-labelledby="auth-title">
        <div className="auth-card">
          <span className="auth-eyebrow">开启你的金融探索</span>
          <h2 id="auth-title">{isRegister ? '创建账号' : '欢迎回来'}</h2>
          <p className="auth-description">{isRegister ? '注册智投顾问，让思路从这里开始。' : '登录智投顾问，开始今天的对话。'}</p>
          {message && <p className="auth-notice" role="status">{message}</p>}
          <form onSubmit={submit} aria-busy={busy}>
            <fieldset disabled={busy} className="auth-fields">
              <label htmlFor="username">用户名</label>
              <input id="username" name="username" autoComplete="username" autoCapitalize="none" spellCheck={false}
                placeholder="请输入用户名" value={username} onChange={(e) => setUsername(e.target.value)} required aria-describedby="username-hint" />
              <small id="username-hint">1–64 个字符，登录时区分大小写</small>
              <label htmlFor="password">密码</label>
              <div className="password-field">
                <input ref={passwordRef} id="password" name="password" type={visible ? 'text' : 'password'}
                  autoComplete={isRegister ? 'new-password' : 'current-password'} placeholder="请输入密码"
                  value={password} onChange={(e) => setPassword(e.target.value)} required aria-describedby={isRegister ? 'password-hint' : undefined} />
                <button className="password-toggle" type="button" onClick={() => setVisible(!visible)} aria-pressed={visible} aria-label={visible ? '隐藏密码' : '显示密码'}>{visible ? '隐藏' : '显示'}</button>
              </div>
              {isRegister && <>
                <small id="password-hint">最多 72 字节，建议混合使用字母、数字和符号</small>
                <label htmlFor="confirmation">确认密码</label>
                <input id="confirmation" name="confirmation" type={visible ? 'text' : 'password'} autoComplete="new-password"
                  placeholder="请再次输入密码" value={confirmation} onChange={(e) => setConfirmation(e.target.value)} required />
              </>}
              {error && <p className="auth-error" role="alert">{error}</p>}
              <button className="auth-submit" type="submit">{busy ? (isRegister ? '正在注册…' : '正在登录…') : (isRegister ? '创建账号' : '登录')}<span aria-hidden="true"> →</span></button>
            </fieldset>
          </form>
          <p className="auth-switch">{isRegister ? '已有账号？' : '还没有账号？'}<button type="button" disabled={busy} onClick={switchMode}>{isRegister ? '去登录' : '立即注册'}</button></p>
          <p className="auth-footnote">登录后将在当前浏览器保持登录，直到凭证过期或主动退出。</p>
        </div>
      </section>
    </main>
  )
}

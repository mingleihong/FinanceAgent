interface Props {
  username: string
  onLogout: () => void
}

export default function Header({ username, onLogout }: Props) {
  return (
    <header className="header">
      <div className="header-inner">
        <div className="brand">
          <span className="brand-mark" aria-hidden="true">
            💹
          </span>
          <div className="brand-text">
            <h1 className="brand-name">智投顾问</h1>
            <p className="brand-tagline">AI 金融建议助手</p>
          </div>
        </div>
        <div className="account-actions">
          <span className="account-name" title={username}>{username}</span>
          <button type="button" className="logout-button" onClick={onLogout}>退出登录</button>
        </div>
      </div>
    </header>
  )
}

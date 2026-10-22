export default function Header() {
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
        <span className="status">
          <span className="status-dot" aria-hidden="true" />
          在线
        </span>
      </div>
    </header>
  )
}

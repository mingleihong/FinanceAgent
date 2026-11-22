import type { SessionItem } from '../api/session'

interface Props {
  sessions: SessionItem[]
  total: number
  page: number
  pageSize: number
  loading: boolean
  error: string
  activeId: string | null
  onPageChange: (page: number) => void
  onRefresh: () => void
  onSelect: (id: string) => void
  onNewChat: () => void
}

function formatTime(value: string): string {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export default function SessionList({
  sessions, total, page, pageSize, loading, error, activeId,
  onPageChange, onRefresh, onSelect, onNewChat,
}: Props) {
  const totalPages = Math.max(1, Math.ceil(total / pageSize))

  return (
    <div className="session-panel">
      <div className="session-panel-head">
        <div>
          <h2 className="session-panel-title">会话列表</h2>
          <span className="session-panel-sub">共 {total} 个会话</span>
        </div>
        <div className="session-panel-actions">
          <button type="button" className="session-icon-button" onClick={onRefresh} disabled={loading} title="刷新">
            刷新
          </button>
          <button type="button" className="session-new-button" onClick={onNewChat} disabled={loading}>
            + 新会话
          </button>
        </div>
      </div>

      {error && <p className="session-error">{error}</p>}

      {loading ? (
        <p className="session-placeholder">加载中…</p>
      ) : sessions.length === 0 ? (
        <p className="session-placeholder">还没有会话，点「+ 新会话」开始聊天吧。</p>
      ) : (
        <ul className="session-list">
          {sessions.map((session) => (
            <li key={session.id}>
              <button
                type="button"
                className={`session-item${session.id === activeId ? ' active' : ''}`}
                onClick={() => onSelect(session.id)}
                aria-current={session.id === activeId ? 'page' : undefined}
              >
                <span className="session-item-title" title={session.title}>{session.title}</span>
                <span className="session-item-time">{formatTime(session.updated_at)}</span>
              </button>
            </li>
          ))}
        </ul>
      )}

      <div className="session-pager">
        <button type="button" onClick={() => onPageChange(page - 1)} disabled={page <= 1 || loading}>
          上一页
        </button>
        <span>{page} / {totalPages}</span>
        <button type="button" onClick={() => onPageChange(page + 1)} disabled={page >= totalPages || loading}>
          下一页
        </button>
      </div>
    </div>
  )
}

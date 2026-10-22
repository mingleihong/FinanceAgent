export default function Thinking() {
  return (
    <div className="message assistant">
      <span className="avatar" aria-hidden="true">
        🤖
      </span>
      <div className="bubble thinking">
        <span className="thinking-dots" aria-label="思考中">
          <span />
          <span />
          <span />
        </span>
      </div>
    </div>
  )
}

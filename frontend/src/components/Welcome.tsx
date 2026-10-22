const SUGGESTIONS = [
  '帮我做一份简单的月度预算规划',
  '什么是定投？适合什么样的人？',
  '我月薪 1 万，怎么开始投资理财？',
  '股票、基金、债券有什么区别？',
]

interface WelcomeProps {
  onPick: (text: string) => void
}

export default function Welcome({ onPick }: WelcomeProps) {
  return (
    <div className="welcome">
      <div className="welcome-mark" aria-hidden="true">
        💹
      </div>
      <h2 className="welcome-title">你好，我是你的金融建议助手</h2>
      <p className="welcome-sub">
        可以帮你做预算规划、解释投资概念、分析理财方案。试着问我一些理财问题吧。
      </p>
      <div className="suggestions">
        {SUGGESTIONS.map((text) => (
          <button
            key={text}
            type="button"
            className="suggestion-chip"
            onClick={() => onPick(text)}
          >
            {text}
          </button>
        ))}
      </div>
    </div>
  )
}

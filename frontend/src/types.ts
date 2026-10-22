/** 后端 /api/chat 的响应结构 */
export interface ChatResponse {
  /** 0 表示成功，非 0 表示失败 */
  code: number
  /** 错误描述，成功时为 "success" */
  message?: string
  /** 助手的最终回答文本，仅成功时存在 */
  data?: string
}

export type Role = 'user' | 'assistant'

export interface Message {
  id: string
  role: Role
  content: string
  /** 标记助手消息出错（如后端不可用） */
  error?: boolean
}

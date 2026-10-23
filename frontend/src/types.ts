export type ChatEvent = 'reasoning_content' | 'content'

export interface Message {
  id: string
  role: 'user' | 'assistant'
  content: string
  reasoning_content?: string
  streaming?: boolean
  error?: string
}

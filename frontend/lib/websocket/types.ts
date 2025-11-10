export type ConnectionState = 'disconnected' | 'connecting' | 'connected' | 'error'

export interface AuthMessage {
  type: 'auth'
  userId: string
  token: string
}

export interface PresenceMessage {
  type: 'presence'
}

export interface ChatMessage {
  type: 'chat'
  recipientId: string
  msg: string
}

export interface IncomingMessage {
  type: 'message'
  data: {
    messageId: string
    senderId: string
    recipientId: string
    content: string
    timestamp: string
    type: string
  }
}

export interface AuthAckMessage {
  type: 'auth_ack'
  status: 'success' | 'failed'
  userId?: string
  error?: string
}

export interface NotificationMessage {
  type: 'notification'
  subType: 'inbox'
  count: number
}

export type WebSocketMessage = AuthMessage | PresenceMessage | ChatMessage | IncomingMessage | AuthAckMessage | NotificationMessage

export interface Message {
  messageId: string
  senderId: string
  recipientId: string
  content: string
  timestamp: Date
  type: string
  direction: 'sent' | 'received'
}


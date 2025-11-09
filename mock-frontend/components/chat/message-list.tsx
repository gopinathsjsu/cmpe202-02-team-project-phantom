'use client'

import { Message } from '@/lib/websocket/types'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { format } from 'date-fns'

interface MessageListProps {
  messages: Message[]
  currentUserId: string
}

export function MessageList({ messages, currentUserId }: MessageListProps) {
  if (messages.length === 0) {
    return (
      <Card className="flex-1">
        <CardContent className="flex items-center justify-center h-full p-8">
          <p className="text-muted-foreground">No messages yet. Start a conversation!</p>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="flex-1 overflow-hidden flex flex-col">
      <CardContent className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((message) => {
          const isSent = message.direction === 'sent'
          const isFromCurrentUser = message.senderId === currentUserId

          return (
            <div
              key={message.messageId}
              className={`flex ${isSent || isFromCurrentUser ? 'justify-end' : 'justify-start'}`}
            >
              <div
                className={`max-w-[70%] rounded-lg p-3 ${
                  isSent || isFromCurrentUser
                    ? 'bg-primary text-primary-foreground'
                    : 'bg-muted text-muted-foreground'
                }`}
              >
                <div className="flex items-center gap-2 mb-1">
                  <Badge variant="outline" className="text-xs">
                    {isSent || isFromCurrentUser ? 'You' : message.senderId}
                  </Badge>
                  {message.direction === 'sent' && (
                    <Badge variant="secondary" className="text-xs">
                      Sent
                    </Badge>
                  )}
                </div>
                <p className="text-sm">{message.content}</p>
                <p className="text-xs opacity-70 mt-1">
                  {format(message.timestamp, 'HH:mm:ss')}
                </p>
              </div>
            </div>
          )
        })}
      </CardContent>
    </Card>
  )
}


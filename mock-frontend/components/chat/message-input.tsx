'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent } from '@/components/ui/card'
import { Send } from 'lucide-react'

interface MessageInputProps {
  onSend: (recipientId: string, content: string) => void
  disabled?: boolean
  recipientId?: string
}

export function MessageInput({ onSend, disabled, recipientId }: MessageInputProps) {
  const [content, setContent] = useState('')
  const [targetRecipientId, setTargetRecipientId] = useState(recipientId || '')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!content.trim() || !targetRecipientId.trim() || disabled) {
      return
    }
    onSend(targetRecipientId.trim(), content.trim())
    setContent('')
  }

  return (
    <Card>
      <CardContent className="p-4">
        {disabled && (
          <p className="text-sm text-muted-foreground mb-2">
            Message sending is disabled. Please ensure WebSocket is connected.
          </p>
        )}
        <form onSubmit={handleSubmit} className="flex gap-2">
          <div className="flex-1 space-y-2">
            <Input
              type="text"
              placeholder="Recipient User ID"
              value={targetRecipientId}
              onChange={(e) => setTargetRecipientId(e.target.value)}
              disabled={disabled}
              className="text-sm"
            />
            <Input
              type="text"
              placeholder="Type a message..."
              value={content}
              onChange={(e) => setContent(e.target.value)}
              disabled={disabled}
            />
          </div>
          <Button type="submit" disabled={disabled || !content.trim() || !targetRecipientId.trim()}>
            <Send className="h-4 w-4" />
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}


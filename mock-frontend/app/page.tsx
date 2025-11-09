'use client'

import { useState, useEffect } from 'react'
import { useAuth } from '@/hooks/use-auth'
import { useWebSocketConnection } from '@/hooks/use-websocket-connection'
import { orchestratorApi } from '@/lib/api/orchestrator'
import { LoginForm } from '@/components/auth/login-form'
import { SignupForm } from '@/components/auth/signup-form'
import { MessageList } from '@/components/chat/message-list'
import { MessageInput } from '@/components/chat/message-input'
import { ConnectionStatus } from '@/components/connection/connection-status'
import { TestControls } from '@/components/testing/test-controls'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { LogOut, Inbox } from 'lucide-react'

export default function HomePage() {
  const [showSignup, setShowSignup] = useState(false)
  const { user, token, refreshToken, isAuthenticated, login, signup, logout } = useAuth()
  
  // Simple WebSocket connection hook - auto-connects when userId and token are available
  const {
    connectionState,
    messages,
    lastHeartbeat,
    connectionError,
    notification,
    connect,
    disconnect,
    sendMessage,
    sendHeartbeat,
    clearMessages,
  } = useWebSocketConnection(user?.user_id || null, token, refreshToken)

  // Fetch undelivered messages at login
  useEffect(() => {
    if (isAuthenticated && token) {
      orchestratorApi
        .getUndeliveredMessages(token, refreshToken)
        .then((result) => {
          if (result.count > 0) {
            console.log(`[HomePage] Found ${result.count} undelivered messages`)
            // Messages will be delivered via WebSocket when user connects
          }
        })
        .catch((error) => {
          console.error('[HomePage] Failed to fetch undelivered messages:', error)
        })
    }
  }, [isAuthenticated, token, refreshToken])

  const messagesSent = messages.filter((m) => m.direction === 'sent').length
  const messagesReceived = messages.filter((m) => m.direction === 'received').length

  if (!isAuthenticated || !user) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background p-4">
        {showSignup ? (
          <SignupForm
            onSignup={async (userName, email, password, phone) => {
              return await signup(userName, email, password, phone)
            }}
            onSwitchToLogin={() => setShowSignup(false)}
          />
        ) : (
          <LoginForm
            onLogin={async (email, password) => {
              return await login(email, password)
            }}
            onSwitchToSignup={() => setShowSignup(true)}
          />
        )}
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background p-4">
      <div className="max-w-7xl mx-auto space-y-4">
        {/* Header */}
        <Card>
          <CardContent className="p-4 flex items-center justify-between">
            <div className="flex items-center gap-4">
              <div>
                <h1 className="text-2xl font-bold">Chat Testing Dashboard</h1>
                <p className="text-sm text-muted-foreground">
                  Logged in as: {user.user_name} ({user.email})
                </p>
                <p className="text-xs text-muted-foreground">User ID: {user.user_id}</p>
              </div>
              {notification && notification.count > 0 && (
                <Badge variant="default" className="flex items-center gap-2">
                  <Inbox className="h-4 w-4" />
                  {notification.count} undelivered message{notification.count !== 1 ? 's' : ''}
                </Badge>
              )}
            </div>
            <Button onClick={logout} variant="outline">
              <LogOut className="h-4 w-4 mr-2" />
              Logout
            </Button>
          </CardContent>
        </Card>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
          {/* Main Chat Area */}
          <div className="lg:col-span-2 space-y-4">
            <MessageList messages={messages} currentUserId={user.user_id} />
            <MessageInput
              onSend={(recipientId, content) => {
                sendMessage(recipientId, content)
              }}
              disabled={connectionState !== 'connected'}
              recipientId=""
            />
          </div>

          {/* Sidebar */}
          <div className="space-y-4">
            <ConnectionStatus
              state={connectionState}
              lastHeartbeat={lastHeartbeat}
              messagesSent={messagesSent}
              messagesReceived={messagesReceived}
              error={connectionError}
            />
            <TestControls
              connectionState={connectionState}
              onConnect={async () => {
                await connect()
              }}
              onDisconnect={disconnect}
              onSendHeartbeat={sendHeartbeat}
              onClearMessages={clearMessages}
            />
            <Card>
              <CardHeader>
                <CardTitle>Testing Instructions</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground space-y-2">
                <p>1. Open multiple browser tabs/windows</p>
                <p>2. Login with different users in each tab</p>
                <p>3. Send messages between users</p>
                <p>4. Test offline scenarios by disconnecting</p>
                <p>5. Check MongoDB for message status</p>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  )
}


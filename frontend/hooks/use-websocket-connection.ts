import { useState, useEffect, useRef, useCallback } from 'react'
import { WebSocketClient } from '@/lib/websocket/client'
import type { ConnectionState, Message, NotificationMessage } from '@/lib/websocket/types'
import { setTokenUpdateCallback } from '@/lib/api/orchestrator'

const EVENTS_SERVER_URL = process.env.NEXT_PUBLIC_EVENTS_SERVER_URL || 'ws://localhost:8001/ws'
const HEARTBEAT_INTERVAL = parseInt(process.env.NEXT_PUBLIC_WS_HEARTBEAT_INTERVAL || '5', 10)

/**
 * Simple WebSocket connection hook with auto-connect.
 * Auto-connects when userId and token are available.
 * 
 * @param userId - User ID for authentication (null if not authenticated)
 * @param token - Auth token for authentication (null if not authenticated)
 * @param refreshToken - Refresh token for token refresh (null if not authenticated)
 */
export function useWebSocketConnection(
  userId: string | null,
  token: string | null,
  refreshToken?: string | null
) {
  const [connectionState, setConnectionState] = useState<ConnectionState>('disconnected')
  const [messages, setMessages] = useState<Message[]>([])
  const [lastHeartbeat, setLastHeartbeat] = useState<Date | null>(null)
  const [connectionError, setConnectionError] = useState<string | null>(null)
  const [notification, setNotification] = useState<NotificationMessage | null>(null)

  const clientRef = useRef<WebSocketClient | null>(null)
  const autoConnectAttemptedRef = useRef(false)
  const tokenUpdateCallbackRef = useRef<((newToken: string, newRefreshToken: string) => void) | null>(null)

  // Set up token update callback for WebSocket client
  useEffect(() => {
    const callback = (newToken: string, newRefreshToken: string) => {
      // Update token in WebSocket client if connected
      if (clientRef.current) {
        clientRef.current['token'] = newToken
        clientRef.current['refreshToken'] = newRefreshToken
      }
      // Also update via API callback
      tokenUpdateCallbackRef.current?.(newToken, newRefreshToken)
    }
    setTokenUpdateCallback(callback)
    tokenUpdateCallbackRef.current = callback
  }, [])

  // Simple check: do we have credentials?
  const hasCredentials = !!userId && !!token

  // Debug logging - only log when credentials actually change or connection state changes significantly
  // This prevents excessive logging on every render
  const prevCredentialsRef = useRef<{ hasCredentials: boolean; userId: string | null; hasToken: boolean } | null>(null)
  useEffect(() => {
    const current = { hasCredentials, userId, hasToken: !!token }
    const prev = prevCredentialsRef.current
    
    // Only log if credentials changed or connection state changed to/from connected
    if (!prev || 
        prev.hasCredentials !== current.hasCredentials ||
        prev.userId !== current.userId ||
        prev.hasToken !== current.hasToken) {
      console.log('[useWebSocketConnection] Credentials changed:', {
        hasCredentials,
        userId,
        hasToken: !!token,
        hasClient: !!clientRef.current,
        connectionState,
        autoConnectAttempted: autoConnectAttemptedRef.current,
      })
      prevCredentialsRef.current = current
    }
  }, [hasCredentials, userId, token, connectionState])

  // Create client when credentials are available
  useEffect(() => {
    if (!hasCredentials) {
      // No credentials - clean up only if client exists
      if (clientRef.current) {
        console.log('[useWebSocketConnection] Cleaning up - no credentials')
        clientRef.current.disconnect()
        clientRef.current = null
      }
      autoConnectAttemptedRef.current = false
      setConnectionState('disconnected')
      setConnectionError(null)
      return
    }

    // We have credentials - create client if needed
    if (!clientRef.current) {
      console.log('[useWebSocketConnection] Creating client', { userId })
      
      const client = new WebSocketClient(
        EVENTS_SERVER_URL,
        HEARTBEAT_INTERVAL,
        {
          onStateChange: (state) => {
            console.log('[useWebSocketConnection] State changed:', state)
            setConnectionState(state)
            if (state === 'connected') {
              setConnectionError(null)
              autoConnectAttemptedRef.current = false
            }
          },
          onMessage: (message) => {
            setMessages((prev) => [...prev, message])
          },
          onError: (error) => {
            console.error('[useWebSocketConnection] Error:', error)
            setConnectionError(error.message)
          },
          onHeartbeat: () => {
            setLastHeartbeat(new Date())
          },
          onNotification: (notification) => {
            console.log('[useWebSocketConnection] Notification received:', notification)
            setNotification(notification)
          },
        },
        tokenUpdateCallbackRef.current || undefined
      )

      clientRef.current = client
      console.log('[useWebSocketConnection] Client created successfully')
    }
  }, [hasCredentials, userId, token])

  // Separate effect for auto-connect to avoid cleanup issues
  useEffect(() => {
    if (!hasCredentials || !clientRef.current) {
      return
    }

    // Auto-connect if disconnected and haven't attempted yet
    if (connectionState === 'disconnected' && !autoConnectAttemptedRef.current) {
      console.log('[useWebSocketConnection] Auto-connecting...', { userId, hasToken: !!token })
      autoConnectAttemptedRef.current = true

      clientRef.current
        .connect(userId!, token!, refreshToken || null)
        .then(() => {
          console.log('[useWebSocketConnection] Auto-connect successful')
        })
        .catch((error) => {
          console.error('[useWebSocketConnection] Auto-connect failed:', error)
          setConnectionError(error instanceof Error ? error.message : 'Auto-connect failed')
          // Allow retry after delay
          setTimeout(() => {
            autoConnectAttemptedRef.current = false
          }, 2000)
        })
    }
  }, [hasCredentials, userId, token, refreshToken, connectionState])

  // Manual connect function
  const connect = useCallback(async () => {
    if (!hasCredentials) {
      const errorMsg = 'Cannot connect: missing userId or token'
      console.warn('[useWebSocketConnection]', errorMsg)
      setConnectionError(errorMsg)
      return { success: false, error: errorMsg }
    }

    if (!clientRef.current) {
      // Create client on-demand
      const client = new WebSocketClient(
        EVENTS_SERVER_URL,
        HEARTBEAT_INTERVAL,
        {
          onStateChange: setConnectionState,
          onMessage: (msg) => setMessages((prev) => [...prev, msg]),
          onError: (err) => setConnectionError(err.message),
          onHeartbeat: () => setLastHeartbeat(new Date()),
          onNotification: (notif) => setNotification(notif),
        },
        tokenUpdateCallbackRef.current || undefined
      )
      clientRef.current = client
    }

    if (connectionState === 'connected' || connectionState === 'connecting') {
      return { success: true }
    }

    console.log('[useWebSocketConnection] Manual connect...', { userId })
    setConnectionError(null)
    autoConnectAttemptedRef.current = true

    try {
      await clientRef.current.connect(userId!, token!, refreshToken || null)
      console.log('[useWebSocketConnection] Manual connect successful')
      return { success: true }
    } catch (error) {
      const errorMsg = error instanceof Error ? error.message : 'Connection failed'
      console.error('[useWebSocketConnection] Manual connect failed:', errorMsg)
      setConnectionError(errorMsg)
      autoConnectAttemptedRef.current = false
      return { success: false, error: errorMsg }
    }
  }, [hasCredentials, userId, token, refreshToken, connectionState])

  const disconnect = useCallback(() => {
    console.log('[useWebSocketConnection] Disconnect')
    if (clientRef.current) {
      clientRef.current.disconnect()
    }
    autoConnectAttemptedRef.current = false
  }, [])

  const sendMessage = useCallback(
    (recipientId: string, content: string) => {
      if (!clientRef.current || connectionState !== 'connected') {
        return { success: false, error: 'Not connected' }
      }

      try {
        clientRef.current.sendChatMessage(recipientId, content)
        const sentMessage: Message = {
          messageId: `temp-${Date.now()}`,
          senderId: userId || '',
          recipientId,
          content,
          timestamp: new Date(),
          type: 'text',
          direction: 'sent',
        }
        setMessages((prev) => [...prev, sentMessage])
        return { success: true }
      } catch (error) {
        return {
          success: false,
          error: error instanceof Error ? error.message : 'Failed to send message',
        }
      }
    },
    [userId, connectionState]
  )

  const sendHeartbeat = useCallback(() => {
    if (!clientRef.current) {
      return { success: false, error: 'Not connected' }
    }
    try {
      clientRef.current.sendPresenceHeartbeat()
      setLastHeartbeat(new Date())
      return { success: true }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to send heartbeat',
      }
    }
  }, [])

  const clearMessages = useCallback(() => {
    setMessages([])
  }, [])

  return {
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
  }
}


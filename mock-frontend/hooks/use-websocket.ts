import { useState, useEffect, useRef, useCallback } from 'react'
import { WebSocketClient } from '@/lib/websocket/client'
import type { ConnectionState, Message } from '@/lib/websocket/types'

const EVENTS_SERVER_URL = process.env.NEXT_PUBLIC_EVENTS_SERVER_URL || 'ws://localhost:8001/ws'
// Default to 5 seconds to match events-server WS_HEARTBEAT_SECONDS
// This should be less than WS_DEAD_SECONDS to keep connection alive
const HEARTBEAT_INTERVAL = parseInt(process.env.NEXT_PUBLIC_WS_HEARTBEAT_INTERVAL || '5', 10)

export function useWebSocket(userId: string | null, token: string | null) {
  const [connectionState, setConnectionState] = useState<ConnectionState>('disconnected')
  const [messages, setMessages] = useState<Message[]>([])
  const [lastHeartbeat, setLastHeartbeat] = useState<Date | null>(null)
  const [connectionError, setConnectionError] = useState<string | null>(null)
  const clientRef = useRef<WebSocketClient | null>(null)
  const isClientReadyRef = useRef(false)

  useEffect(() => {
    if (!userId || !token) {
      clientRef.current = null
      isClientReadyRef.current = false
      return
    }

    // Create client synchronously when userId/token are available
    const client = new WebSocketClient(EVENTS_SERVER_URL, HEARTBEAT_INTERVAL, {
      onStateChange: (state) => {
        setConnectionState(state)
        // Clear error when state changes to connected
        if (state === 'connected') {
          setConnectionError(null)
        }
      },
      onMessage: (message) => {
        setMessages((prev) => [...prev, message])
      },
      onError: (error) => {
        console.error('[useWebSocket] WebSocket error:', error)
        setConnectionError(error.message)
      },
      onHeartbeat: () => {
        setLastHeartbeat(new Date())
      },
    })

    clientRef.current = client
    isClientReadyRef.current = true
    console.log('[useWebSocket] WebSocket client created and ready', { userId })

    return () => {
      client.disconnect()
      clientRef.current = null
      isClientReadyRef.current = false
    }
  }, [userId, token])

  const connect = useCallback(async () => {
    // Wait a bit if client is not ready yet (should be very rare)
    if (!isClientReadyRef.current || !clientRef.current) {
      console.warn('[useWebSocket] Client not ready yet, waiting...', {
        isReady: isClientReadyRef.current,
        hasClient: !!clientRef.current,
        userId,
        hasToken: !!token,
      })
      // Wait a short time for client to be created
      await new Promise((resolve) => setTimeout(resolve, 100))
      if (!isClientReadyRef.current || !clientRef.current) {
        const errorMsg = 'WebSocket client not initialized. Please try again.'
        console.error('[useWebSocket]', errorMsg)
        setConnectionError(errorMsg)
        return { success: false, error: errorMsg }
      }
    }

    if (!userId || !token) {
      const errorMsg = 'Missing userId or token'
      console.warn('[useWebSocket] Cannot connect:', errorMsg)
      setConnectionError(errorMsg)
      return { success: false, error: errorMsg }
    }

    console.log('[useWebSocket] Attempting to connect...', { userId })
    setConnectionError(null)
    try {
      await clientRef.current.connect(userId, token)
      console.log('[useWebSocket] Connection successful')
      setConnectionError(null)
      return { success: true }
    } catch (error) {
      const errorMsg = error instanceof Error ? error.message : 'Connection failed'
      console.error('[useWebSocket] Connection failed:', errorMsg)
      setConnectionError(errorMsg)
      return { success: false, error: errorMsg }
    }
  }, [userId, token])

  const disconnect = useCallback(() => {
    clientRef.current?.disconnect()
  }, [])

  const sendMessage = useCallback((recipientId: string, content: string) => {
    if (!clientRef.current) {
      return { success: false, error: 'WebSocket not initialized' }
    }

    try {
      clientRef.current.sendChatMessage(recipientId, content)
      
      // Add sent message to local state
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
      return { success: false, error: error instanceof Error ? error.message : 'Failed to send message' }
    }
  }, [userId])

  const sendHeartbeat = useCallback(() => {
    if (!clientRef.current) {
      return { success: false, error: 'WebSocket not initialized' }
    }

    try {
      clientRef.current.sendPresenceHeartbeat()
      setLastHeartbeat(new Date())
      return { success: true }
    } catch (error) {
      return { success: false, error: error instanceof Error ? error.message : 'Failed to send heartbeat' }
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
    connect,
    disconnect,
    sendMessage,
    sendHeartbeat,
    clearMessages,
  }
}


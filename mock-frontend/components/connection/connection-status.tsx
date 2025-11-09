'use client'

import { ConnectionState } from '@/lib/websocket/types'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Wifi, WifiOff, Loader2, AlertCircle } from 'lucide-react'

interface ConnectionStatusProps {
  state: ConnectionState
  lastHeartbeat: Date | null
  messagesSent: number
  messagesReceived: number
  error?: string | null
}

export function ConnectionStatus({
  state,
  lastHeartbeat,
  messagesSent,
  messagesReceived,
  error,
}: ConnectionStatusProps) {
  const getStateIcon = () => {
    switch (state) {
      case 'connected':
        return <Wifi className="h-4 w-4 text-green-500" />
      case 'connecting':
        return <Loader2 className="h-4 w-4 text-yellow-500 animate-spin" />
      case 'error':
        return <AlertCircle className="h-4 w-4 text-red-500" />
      default:
        return <WifiOff className="h-4 w-4 text-gray-500" />
    }
  }

  const getStateColor = () => {
    switch (state) {
      case 'connected':
        return 'bg-green-500/10 text-green-500 border-green-500/20'
      case 'connecting':
        return 'bg-yellow-500/10 text-yellow-500 border-yellow-500/20'
      case 'error':
        return 'bg-red-500/10 text-red-500 border-red-500/20'
      default:
        return 'bg-gray-500/10 text-gray-500 border-gray-500/20'
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          Connection Status
          {getStateIcon()}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">State:</span>
          <Badge variant="outline" className={getStateColor()}>
            {state.toUpperCase()}
          </Badge>
        </div>
        {lastHeartbeat && (
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Last Heartbeat:</span>
            <span className="text-sm">
              {lastHeartbeat.toLocaleTimeString()}
            </span>
          </div>
        )}
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">Messages Sent:</span>
          <Badge variant="secondary">{messagesSent}</Badge>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">Messages Received:</span>
          <Badge variant="secondary">{messagesReceived}</Badge>
        </div>
        {error && (
          <div className="mt-3 p-2 bg-destructive/10 border border-destructive/20 rounded-md">
            <p className="text-xs text-destructive font-medium">Connection Error:</p>
            <p className="text-xs text-destructive/80">{error}</p>
          </div>
        )}
      </CardContent>
    </Card>
  )
}


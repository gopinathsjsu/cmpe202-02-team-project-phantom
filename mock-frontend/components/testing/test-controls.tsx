'use client'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ConnectionState } from '@/lib/websocket/types'
import { Wifi, WifiOff, Heart, Trash2 } from 'lucide-react'

interface TestControlsProps {
  connectionState: ConnectionState
  onConnect: () => void
  onDisconnect: () => void
  onSendHeartbeat: () => void
  onClearMessages: () => void
}

export function TestControls({
  connectionState,
  onConnect,
  onDisconnect,
  onSendHeartbeat,
  onClearMessages,
}: TestControlsProps) {
  const isConnected = connectionState === 'connected'
  const isConnecting = connectionState === 'connecting'

  return (
    <Card>
      <CardHeader>
        <CardTitle>Test Controls</CardTitle>
      </CardHeader>
      <CardContent className="space-y-2">
        <div className="flex gap-2">
          {!isConnected && !isConnecting && (
            <Button onClick={onConnect} className="flex-1">
              <Wifi className="h-4 w-4 mr-2" />
              Connect
            </Button>
          )}
          {isConnected && (
            <Button onClick={onDisconnect} variant="destructive" className="flex-1">
              <WifiOff className="h-4 w-4 mr-2" />
              Disconnect
            </Button>
          )}
        </div>
        <Button
          onClick={onSendHeartbeat}
          variant="outline"
          className="w-full"
          disabled={!isConnected}
          title="Manual heartbeat (automatic heartbeats run every 5 seconds when connected)"
        >
          <Heart className="h-4 w-4 mr-2" />
          Send Heartbeat (Manual)
        </Button>
        <p className="text-xs text-muted-foreground">
          Automatic heartbeats are sent every 5 seconds when connected (configurable via NEXT_PUBLIC_WS_HEARTBEAT_INTERVAL)
        </p>
        <Button
          onClick={onClearMessages}
          variant="outline"
          className="w-full"
        >
          <Trash2 className="h-4 w-4 mr-2" />
          Clear Messages
        </Button>
      </CardContent>
    </Card>
  )
}


# Mock Frontend - Chat Testing Application

A Next.js/TypeScript application for testing the events-server and chat-consumer systems. This application provides a UI for testing authentication, WebSocket connections, message delivery, online/offline behavior, and reconnection scenarios.

## Features

- **Authentication**: Sign up and login using the orchestrator API
- **WebSocket Integration**: Connect to events-server and manage WebSocket lifecycle
- **Message Testing**: Send and receive messages with visual feedback
- **Presence Testing**: Visual indicators for online/offline status
- **Multi-User Testing**: Support multiple concurrent sessions (open multiple browser tabs)
- **Test Controls**: Easy controls for testing different scenarios
- **Status Monitoring**: Real-time connection and message status

## Setup

1. Install dependencies:
```bash
npm install
# or
pnpm install
```

2. Create a `.env.local` file (copy from `env.example`):
```bash
cp env.example .env.local
```

3. Configure environment variables in `.env.local`:
```env
NEXT_PUBLIC_ORCHESTRATOR_URL=http://localhost:8080
NEXT_PUBLIC_EVENTS_SERVER_URL=ws://localhost:8001/ws
NEXT_PUBLIC_WS_HEARTBEAT_INTERVAL=30
```

   **Note**: Update the ports to match your backend services:
   - `NEXT_PUBLIC_ORCHESTRATOR_URL`: Should match orchestrator's PORT (default: 8080)
   - `NEXT_PUBLIC_EVENTS_SERVER_URL`: Should match events-server's PORT (check events-server config)
   
   See `ENV_SETUP.md` for detailed information about all environment variables.

4. Start the development server:
```bash
npm run dev
# or
pnpm dev
```

5. Open [http://localhost:3000](http://localhost:3000) in your browser

## Usage

### Testing Workflow

1. **Sign Up/Login**: Create an account or login with existing credentials
2. **Auto-Connect**: The WebSocket will automatically connect after authentication
3. **Send Messages**: 
   - Enter a recipient User ID
   - Type a message and send
4. **Multi-User Testing**:
   - Open multiple browser tabs/windows
   - Login with different users in each tab
   - Send messages between users
5. **Test Scenarios**:
   - **Online Delivery**: Send message when recipient is online
   - **Offline Handling**: Disconnect recipient, send message, check MongoDB for UNDELIVERED status
   - **Reconnection**: Disconnect and reconnect, verify message delivery
   - **Presence Heartbeat**: Use "Send Heartbeat" button to manually trigger presence
   - **Connection Timeout**: Test WSDeadSeconds timeout behavior

### Testing Scenarios

#### Online Message Delivery
1. Open two browser tabs
2. Login with User A in tab 1, User B in tab 2
3. From User A, send a message to User B
4. Message should appear in User B's tab immediately
5. Check MongoDB - message should have status "DELIVERED"

#### Offline Message Handling
1. Open two browser tabs
2. Login with User A in tab 1, User B in tab 2
3. Disconnect User B's WebSocket (click "Disconnect" in tab 2)
4. From User A, send a message to User B
5. Check MongoDB - message should have status "UNDELIVERED"
6. Reconnect User B
7. Message should be delivered when User B comes online

#### Reconnection Testing
1. Connect to WebSocket
2. Send some messages
3. Disconnect WebSocket
4. Reconnect WebSocket
5. Verify connection state and message delivery

#### Presence Heartbeat
1. Connect to WebSocket
2. Observe automatic heartbeats (every 30 seconds by default)
3. Use "Send Heartbeat" button to manually trigger presence
4. Check Redis for presence keys: `presence:{userId}`

## Architecture

### Message Flow

1. **Client → Events-Server**: WebSocket connection with auth message
2. **Events-Server → RabbitMQ**: Chat messages are queued
3. **Chat-Consumer**: Consumes from RabbitMQ, saves to MongoDB
4. **Presence Check**: Chat-consumer checks if recipient is online
5. **If Online**: Redis pub/sub → Events-Server → WebSocket → Client
6. **If Offline**: Message marked UNDELIVERED in MongoDB

### Components

- **Authentication**: Login/Signup forms with orchestrator API integration
- **WebSocket Client**: Manages connection, auth, presence, and messages
- **Chat Interface**: Message list and input components
- **Connection Status**: Real-time connection state and metrics
- **Test Controls**: Buttons for testing scenarios

## Environment Variables

- `NEXT_PUBLIC_ORCHESTRATOR_URL`: Base URL for orchestrator API (default: http://localhost:8080)
- `NEXT_PUBLIC_EVENTS_SERVER_URL`: WebSocket URL for events-server (default: ws://localhost:8001/ws)
- `NEXT_PUBLIC_WS_HEARTBEAT_INTERVAL`: Presence heartbeat interval in seconds (default: 30)

## Dependencies

- Next.js 16
- React 19
- TypeScript 5
- Tailwind CSS
- date-fns (for date formatting)
- lucide-react (for icons)

## Notes

- Messages are stored in MongoDB by chat-consumer
- Presence is tracked in Redis by events-server
- Message status can be checked in MongoDB: SENT, DELIVERED, UNDELIVERED
- For multi-user testing, open multiple browser tabs/windows
- Each tab represents a different user session


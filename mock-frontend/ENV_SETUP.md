# Environment Setup Guide

This document explains the environment variables needed for mock-frontend and how they relate to the backend services.

## Mock-Frontend Environment Variables

Create a `.env.local` file in the `mock-frontend` directory with the following variables:

```env
# Orchestrator API URL (for authentication)
NEXT_PUBLIC_ORCHESTRATOR_URL=http://localhost:8080

# Events Server WebSocket URL
NEXT_PUBLIC_EVENTS_SERVER_URL=ws://localhost:8001/ws

# WebSocket Heartbeat Interval (seconds)
NEXT_PUBLIC_WS_HEARTBEAT_INTERVAL=30
```

## Backend Service Environment Variables

### Orchestrator

The orchestrator service requires these environment variables:

```env
# Server port (defaults to 8080 if not set)
PORT=8080

# Database connection
DATABASE_URL=postgres://user:password@host:port/dbname?sslmode=disable

# Optional: MongoDB for chat messages
CHAT_MONGO_URI=mongodb://localhost:27017/chatdb

# Optional: RabbitMQ for message queue
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_QUEUE_NAME=chat_messages

# Optional: Listing service integration
LISTING_SERVICE_URL=http://localhost:8082
LISTING_SERVICE_SHARED_SECRET=your-secret
```

**Location**: `orchestrator/cmd/main.go`

### Events-Server

The events-server requires these environment variables:

```env
# WebSocket server port
PORT=:8001

# Orchestrator URL for auth verification
ORCH_BASE_URL=http://localhost:8080

# Redis configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# WebSocket configuration
WS_HEARTBEAT_SECONDS=30
WS_DEAD_SECONDS=60
PRESENCE_TTL_SECONDS=60

# RabbitMQ configuration
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_QUEUE_NAME=chat_messages
```

**Location**: `events-server/internal/config/config.go`

**Note**: The `PORT` format includes the colon (`:8081`), but when constructing the WebSocket URL in mock-frontend, use `ws://localhost:8081/ws` (without the colon in the URL).

### Chat-Consumer

The chat-consumer requires these environment variables:

```env
# RabbitMQ configuration
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_QUEUE_NAME=chat_messages

# Redis configuration (for presence checking)
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# MongoDB configuration
MONGO_URI=mongodb://localhost:27017/chatdb
```

**Location**: `chat-consumer/internal/config/config.go`

## Quick Setup

1. **Copy the example file**:
   ```bash
   cp env.example .env.local
   ```

2. **Update the values** in `.env.local` to match your backend service ports:
   - If orchestrator runs on a different port, update `NEXT_PUBLIC_ORCHESTRATOR_URL`
   - If events-server runs on a different port, update `NEXT_PUBLIC_EVENTS_SERVER_URL`
   - Adjust heartbeat interval if needed

3. **Verify backend services are running**:
   - Orchestrator should be accessible at the URL specified in `NEXT_PUBLIC_ORCHESTRATOR_URL`
   - Events-server WebSocket should be accessible at the URL specified in `NEXT_PUBLIC_EVENTS_SERVER_URL`

## Port Reference

Based on the codebase:
- **Orchestrator**: Default port `8080` (see `orchestrator/cmd/main.go`)
- **Events-Server**: Check the `PORT` environment variable (commonly `8001`)
- **Chat-Consumer**: No HTTP port (consumes from RabbitMQ only)

## Testing the Configuration

1. Start all backend services (orchestrator, events-server, chat-consumer)
2. Verify orchestrator is accessible: `curl http://localhost:8080/health`
3. Start mock-frontend: `npm run dev`
4. Try logging in - if authentication works, the orchestrator URL is correct
5. After login, check if WebSocket connects - if it does, the events-server URL is correct


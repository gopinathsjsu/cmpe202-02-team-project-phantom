# Environment Variables Setup Guide

This document explains all environment variables needed for the Docker Compose deployment.

## .env File Locations

### Primary .env File
**Location**: `.env` (root directory)  
**Path**: `/Users/danlam/Desktop/School/Personal/CMPE202Project/Actual/cmpe202-02-team-project-phantom/.env`

This is the main environment file used by `docker-compose.yml`. All services read from this file.

### Service-Specific .env Files (Optional)
These are optional and provide service-specific overrides:
- `./orchestrator/.env` - Orchestrator service overrides
- `./listing-service/.env` - Listing service overrides
- `./events-server/.env` - Events server overrides
- `./chat-consumer/.env` - Chat consumer overrides

**Note**: Service-specific .env files are merged with the root .env file by docker-compose.

## Quick Start

1. Copy the example file:
   ```bash
   cp .env.example .env
   ```

2. Fill in all required variables (see list below)

3. Start services:
   ```bash
   docker-compose up -d
   ```

## Required Environment Variables

### 1. Database Configuration
```env
# PostgreSQL/Neon Database (REQUIRED)
DATABASE_URL=postgres://user:password@your-neon-host.neon.tech/dbname?sslmode=require
```

### 2. Authentication Secrets (REQUIRED)
```env
# JWT Access Token Secret - Generate with: openssl rand -hex 32
JWT_TOKEN_SECRET=your-secure-random-string-here

# JWT Refresh Token Secret - Generate with: openssl rand -hex 32
JWT_REFRESH_SECRET=your-secure-random-string-here
```

### 3. MongoDB Atlas (REQUIRED for chat-consumer)
```env
# MongoDB Atlas connection string (REQUIRED for chat-consumer)
MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/chatdb?retryWrites=true&w=majority

# MongoDB Atlas connection string (OPTIONAL for orchestrator - chat features won't work without it)
CHAT_MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/chatdb?retryWrites=true&w=majority
```

### 4. Azure Blob Storage (REQUIRED for listing-service)
```env
# Azure Storage Account URL
AZURE_ACCOUNT_URL=https://yourstorageaccount.blob.core.windows.net

# Azure Storage Account Name
AZURE_ACCOUNTNAME=yourstorageaccount

# Azure Storage Account Key
AZURE_ACCOUNTKEY=your-azure-storage-account-key-here

# Azure Blob Container Name
AZURE_CONTAINERNAME=your-container-name
```

### 5. Google Gemini API (REQUIRED for listing-service AI search)
```env
# Get from: https://makersuite.google.com/app/apikey
GOOGLE_API_KEY=your-google-gemini-api-key-here
```

### 6. Service Communication Secrets
```env
# Shared secret between orchestrator and listing-service
# These should match!
LISTING_SERVICE_SHARED_SECRET=your-shared-secret-here
ORCH_REQUEST_ID=your-shared-secret-here
```

### 7. Deployment Configuration
```env
# Your AWS public IP or domain (for frontend browser connections)
EXTERNAL_HOST=3.17.68.221
```

## Optional Environment Variables (with defaults)

### RabbitMQ Configuration
```env
# RabbitMQ connection (default: amqp://guest:guest@rabbitmq:5672/)
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/

# RabbitMQ queue name (default: chat_messages)
RABBITMQ_QUEUE_NAME=chat_messages
```

### Redis Configuration
```env
# Redis address (default: redis:6379)
REDIS_ADDR=redis:6379

# Redis password (default: empty)
REDIS_PASSWORD=

# Redis database number (default: 0)
REDIS_DB=0
```

### Service Ports
```env
# Orchestrator port (default: 8080)
ORCHESTRATOR_PORT=8080
PORT=8080

# Listing service port (default: 8080)
LISTING_PORT=8080
```

### CORS Configuration
```env
# CORS allowed origins (comma-separated)
# Default: http://localhost:3000,http://localhost:3001
# Production: http://3.17.68.221:3000,https://yourdomain.com
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001,http://3.17.68.221:3000
```

### Events Server Configuration
```env
# Presence TTL in seconds (default: 60)
PRESENCE_TTL_SECONDS=60

# Orchestrator base URL (default: http://orchestrator:8080)
ORCH_BASE_URL=http://orchestrator:8080
```

### Frontend Configuration
Frontend URLs are automatically constructed from `EXTERNAL_HOST` and ports:
- `NEXT_PUBLIC_ORCHESTRATOR_URL`: `http://${EXTERNAL_HOST}:${ORCHESTRATOR_PORT}`
- `NEXT_PUBLIC_EVENTS_SERVER_URL`: `ws://${EXTERNAL_HOST}:8001/ws`

To override, modify `docker-compose.yml` frontend service environment section.

## Service-Specific Variable Breakdown

### Orchestrator Service
**Location**: `./orchestrator/.env` (optional override)

Required:
- `DATABASE_URL`
- `JWT_TOKEN_SECRET`
- `JWT_REFRESH_SECRET`

Optional:
- `CHAT_MONGO_URI` (for chat features)
- `RABBITMQ_URL`
- `RABBITMQ_QUEUE_NAME`
- `LISTING_SERVICE_URL`
- `LISTING_SERVICE_SHARED_SECRET`
- `CORS_ALLOWED_ORIGINS`
- `PORT`

### Listing Service
**Location**: `./listing-service/.env` (optional override)

Required:
- `DATABASE_URL`
- `GOOGLE_API_KEY`
- `AZURE_ACCOUNT_URL`
- `AZURE_ACCOUNTNAME`
- `AZURE_ACCOUNTKEY`
- `AZURE_CONTAINERNAME`
- `ORCH_REQUEST_ID`

Optional:
- `LISTING_PORT`

### Events Server
**Location**: `./events-server/.env` (optional override)

Required:
- `PORT`
- `ORCH_BASE_URL`
- `REDIS_ADDR`
- `REDIS_DB`
- `RABBITMQ_URL`
- `RABBITMQ_QUEUE_NAME`

Optional:
- `REDIS_PASSWORD`
- `PRESENCE_TTL_SECONDS`
- `CORS_ALLOWED_ORIGINS`

### Chat Consumer
**Location**: `./chat-consumer/.env` (optional override)

Required:
- `MONGO_URI`
- `RABBITMQ_URL`
- `RABBITMQ_QUEUE_NAME`
- `REDIS_ADDR`
- `REDIS_DB`

Optional:
- `REDIS_PASSWORD`

## Generating Secure Secrets

### JWT Secrets
```bash
# Generate JWT_TOKEN_SECRET
openssl rand -hex 32

# Generate JWT_REFRESH_SECRET
openssl rand -hex 32
```

### Service Communication Secret
```bash
# Generate shared secret for orchestrator <-> listing-service
openssl rand -hex 32
```

## Production Checklist

Before deploying to production, ensure:

- [ ] All required variables are set
- [ ] `EXTERNAL_HOST` is set to your AWS public IP or domain
- [ ] `CORS_ALLOWED_ORIGINS` includes your production frontend URL
- [ ] JWT secrets are strong and unique
- [ ] MongoDB Atlas allows connections from your AWS IPs
- [ ] Database connection string uses SSL (`sslmode=require`)
- [ ] Azure Storage account is accessible from your deployment
- [ ] All secrets are stored securely (not committed to git)

## Troubleshooting

### Service can't connect to database
- Check `DATABASE_URL` format
- Ensure database is accessible from Docker network
- Verify SSL settings if using Neon/cloud database

### Chat features not working
- Ensure `CHAT_MONGO_URI` is set in orchestrator
- Ensure `MONGO_URI` is set in chat-consumer
- Verify MongoDB Atlas network access allows your AWS IPs

### Frontend can't connect to backend
- Check `EXTERNAL_HOST` matches your AWS public IP
- Verify `CORS_ALLOWED_ORIGINS` includes your frontend URL
- Check that ports are properly exposed in docker-compose.yml

### Media uploads failing
- Verify Azure credentials are correct
- Check container name exists in Azure Storage
- Ensure container permissions allow uploads


# CMPE 202 – Marketplace Platform

A full-stack microservices-based marketplace application featuring real-time messaging, AI-powered search, media uploads, admin moderation, and a responsive Next.js frontend.  
The system is fully containerized with Docker Compose and deployed to AWS EC2 with Neon Postgres and Azure Blob Storage.

---

# Features Overview

- Next.js 14 frontend with dynamic listings, filters, saved listings, and admin tools  
- Multiple Go microservices for Listings, Users, Chat, WebSockets, and Orchestrator  
- Real-time chat using WebSockets and background consumer for message persistence  
- AI search assistant for natural-language queries  
- Secure media upload through Azure Blob Storage using SAS URLs  
- Neon Postgres database with SSL enforcement  
- Full Docker Compose deployment across all services  
- Admin dashboard for flagged listings, user moderation, and analytics  

---

# Prerequisites

### Required
- Docker & Docker Compose
- Go 1.22+
- Node.js 18+
- Neon Postgres (or local Postgres)
- Azure Storage Account (for media)
- AWS account (for production)

### Optional
- Domain + SSL  
- Reverse proxy (Nginx or ALB)

---

# System Architecture

```
root/
├── orchestrator/          # API gateway / request router
├── listing-service/       # Listing CRUD, search, analytics
├── user-service/          # Authentication, profiles, admin actions
├── events-server/         # WebSocket real-time messaging server
├── chat-consumer/         # Queue consumer for saving chat messages
├── frontend/              # Next.js client application
├── http-lib/              # Shared Go helper library
└── docker-compose.yml     # Compose file connecting all services
```

### Core Technologies
- **Frontend:** Next.js, React, Tailwind  
- **Backend:** Go microservices (Mux/Fiber)  
- **Database:** Neon Serverless Postgres  
- **Media Storage:** Azure Blob Storage  
- **Real-Time:** WebSocket events server  
- **Deployment:** Docker Compose on AWS EC2  

---

# Quick Start (Local Development)

### 1. Clone Repository
```bash
git clone https://github.com/your-org/marketplace.git
cd marketplace
```

### 2. Create `.env` Files  
Each service includes `.env.example`. Copy them:

```bash
cp listing-service/.env.example listing-service/.env
cp user-service/.env.example user-service/.env
cp orchestrator/.env.example orchestrator/.env
cp frontend/.env.example frontend/.env
```

### 3. Start Entire System
```bash
docker compose up --build
```

Service Ports:

| Service | Port |
|---------|------|
| Frontend | 3000 |
| Orchestrator | 8080 |
| Listing Service | 8081 |
| User Service | 8082 |
| WebSocket Server | 9090 |

---

# Development by Service

### Frontend (Next.js)
```bash
cd frontend
npm install
npm run dev
```

### Listing Service
```bash
cd listing-service
go mod tidy
go run ./cmd
```

### Orchestrator
```bash
cd orchestrator
go run ./cmd
```

---

# API Endpoints (Overview)

### Listings
- GET /api/listings
- GET /api/listings/{id}
- POST /api/listings/create
- PATCH /api/listings/update/{id}
- DELETE /api/listings/delete/{id}
- POST /api/listings/upload
- POST /api/listings/add-media-url/{id}

### Users
- Login / Register
- Edit profile  
- Get personal listings  
- Admin user management  

### Chat & Messaging
- WebSocket: ws://host:9090/ws
- AI search endpoint: /api/listings/chatsearch
- Create conversations  
- Fetch conversation history  

### Admin
- Flagged listings  
- Combined reports  
- Analytics  

---

# Frontend Feature Set

### Listings
- Create, edit, delete listings  
- Media upload UI  
- Saved listings  
- Price filters with debounced inputs  
- Category navigation  
- Homepage category counts & featured items  

### User
- Login / registration  
- Profile edit  
- Personal listings  

### Chat
- Real-time messaging  
- Conversation list  
- Sequential AI chat assistant  

### Admin
- User management  
- Flagged listing review  
- Listing moderation panel  
- Analytics dashboard  

---

# Deployment (AWS)

### Deployment Steps
1. Launch Ubuntu EC2 instance  
2. Install Docker & Docker Compose  
3. Pull repository  
4. Configure production `.env` files  
5. Run:
```bash
docker compose up --build -d
```
6. Create AWS Application Load Balancer  
7. Route domain → ALB → EC2:8080  
8. Enable HTTPS

---

# Environment Variables (Example)

### Database
```
DATABASE_URL=postgres://user:pass@neon-host/db?sslmode=require
```

### Azure Blob Storage
```
AZURE_STORAGE_ACCOUNT=
AZURE_STORAGE_CONTAINER=
AZURE_SAS_TOKEN=
```

### Authentication / Services
```
JWT_SECRET=
REFRESH_SECRET=
LISTING_SERVICE_URL=http://listing-service:8081
USER_SERVICE_URL=http://user-service:8082
EVENTS_WS_URL=ws://events-server:9090/ws
```

---

# Technology Stack

### Frontend
- Next.js  
- React  
- Tailwind CSS  
- Axios  

### Backend
- Go (Golang)  
- Fiber / Mux  
- Neon Postgres  
- WebSockets  
- Azure Blob Upload  

### Infrastructure
- Docker  
- AWS EC2  
- Application Load Balancer  

---

# UML Diagrams

![Component Diagram](Project-Journal/uml-diagrams/Component%20Diagram%20-%20Campus%20Marketplace.png)
![Component Diagram](Project-Journal/uml-diagrams/Deployment%20Diagram%20-%20Campus%20Marketplace.png)

---

# Security Features
- JWT authentication with refresh tokens  
- SAS-based secure media uploads  
- SSL/TLS enforced DB connections  
- Role-based access control (User/Admin)  
- Request validation and rate limiting  
- CORS and secure headers  

---

# Project Structure

```
marketplace/
├── orchestrator/
├── listing-service/
├── user-service/
├── events-server/
├── chat-consumer/
├── frontend/
├── http-lib/
└── docker-compose.yml
```

---

# Credits / Contributors
- **Dan Lam** – Frontend architecture, listings, admin dashboard, reporting, homepage, UX  
- **Kunal Singh** – Backend architecture, orchestrator, chat infrastructure, user management, Docker  
- **Nikhil Raj Singh** – Frontend chat integration, AI search UI, navigation, real-time UX  

---

# End of README

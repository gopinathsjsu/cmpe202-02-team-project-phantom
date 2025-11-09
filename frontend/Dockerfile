# --- base ---
FROM node:20-bullseye AS base
WORKDIR /app
COPY package*.json ./

# --- deps ---
FROM base AS deps
RUN npm ci --legacy-peer-deps

# --- builder (prod) ---
FROM base AS builder
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# --- dev runtime ---
FROM base AS dev
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ENV NODE_ENV=development
EXPOSE 3000
CMD ["npm","run","dev"]
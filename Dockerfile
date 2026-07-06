# ── Build stage ──────────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS backend-build
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

FROM node:24-alpine AS frontend-build
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:3.21
# rtl-sdr provides rtl_fm; multimon-ng decodes SAME/EAS from the audio stream
RUN apk add --no-cache ca-certificates tzdata rtl-sdr multimon-ng
WORKDIR /app

COPY --from=backend-build /app/server ./server
COPY --from=frontend-build /app/frontend/dist ./frontend/dist

EXPOSE 8080
CMD ["./server"]

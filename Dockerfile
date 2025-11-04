# ---------- Stage 1: Build ----------
FROM golang:1.25.1-alpine AS builder

# Install git (needed for go mod download in private repos)
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binaries
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o didlydoodash-api ./cmd/api

# ---------- Stage 2: Runtime ----------
FROM alpine:3.20

# Create non-root user
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy binaries only
COPY --from=builder /app/didlydoodash-api ./

# Use non-root user
USER appuser

# Expose API port
EXPOSE 3000

# Start the API binary directly
ENTRYPOINT ["./didlydoodash-api"]

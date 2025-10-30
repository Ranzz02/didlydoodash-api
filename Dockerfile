# Stage 1: Build the Go binaries using the official Golang image
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod and go.sum to install dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source files and certificate
COPY . .

# Build the Go binaries for the migrate and api commands
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./didlydoodash-migrate ./cmd/migrate && \
    CGO_ENABLED=0 GOOS=linux go build -v -o ./didlydoodash-api ./cmd/api

# Stage 2: Create a lightweight image using Alpine Linux
FROM alpine:latest

WORKDIR /api

# Copy the built binaries and the startup script from the builder stage
COPY --from=builder /app/didlydoodash-api /app/didlydoodash-migrate /app/start.sh ./

# Ensure the script has unix line endings
RUN sed -i 's/\r$//' start.sh

# Give execution permissions to the startup script
RUN chmod +x start.sh

# Expose the port the API will run on
EXPOSE 3000

# Set the startup script as the container's entry point
ENTRYPOINT ["sh","./start.sh"]
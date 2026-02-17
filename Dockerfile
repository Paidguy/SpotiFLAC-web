# Stage 1: Build frontend
FROM node:20-alpine AS frontend
WORKDIR /app/frontend

# Copy frontend package files
COPY frontend/package.json frontend/pnpm-lock.yaml* ./

# Install pnpm and dependencies
RUN npm install -g pnpm && \
    pnpm install --frozen-lockfile || pnpm install

# Copy frontend source code
COPY frontend/ ./

# Build frontend
RUN pnpm run build

# Stage 2: Build Go binary
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy Go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy built frontend from previous stage
COPY --from=frontend /app/frontend/dist ./frontend/dist

# Build the Go binary
RUN go build -o spotiflac .

# Stage 3: Runtime
FROM alpine:3.19
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ffmpeg ca-certificates

# Copy binary from builder
COPY --from=builder /app/spotiflac .

# Create volume mount points
VOLUME ["/downloads", "/data"]

# Environment variables
ENV DOWNLOAD_PATH=/downloads
ENV DATA_DIR=/data
ENV PORT=8080

# Expose port
EXPOSE 8080

# Run the server
CMD ["./spotiflac"]

# Build stage
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

# Set build arguments
ARG VERSION=dev
ARG BUILD_TIME=unknown
ARG GIT_COMMIT=unknown
ARG TARGETPLATFORM
ARG BUILDPLATFORM

# Install build dependencies
RUN apk add --no-cache make git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN GOOS=$(echo $TARGETPLATFORM | cut -d/ -f1) \
    GOARCH=$(echo $TARGETPLATFORM | cut -d/ -f2) \
    go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT" \
    -o /app/dseq

# Final stage
FROM --platform=$TARGETPLATFORM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/dseq /app/dseq

# Set entrypoint
ENTRYPOINT ["/app/dseq"] 